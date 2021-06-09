package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"k8s.io/client-go/tools/clientcmd"

	"context"

	"reflect"

	"github.com/acornsoftlab/dashboard/terminal/backend/localcommand"
	cache "github.com/acornsoftlab/dashboard/terminal/cache/token"
	"github.com/acornsoftlab/dashboard/terminal/pkg/config"
	"github.com/acornsoftlab/dashboard/terminal/pkg/randomstring"
	"github.com/acornsoftlab/dashboard/terminal/server"
	"github.com/acornsoftlab/dashboard/terminal/utils"
	"github.com/acornsoftlab/dashboard/terminal/webtty"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

//terminal shell request info
type termRequest struct {
	inclustermode string
	kubeconfig    string
	kubetoken     string
	cluster       string
	namespace     string
	pod           string
	container     string
	termtype      string
}

type Response struct {
	Success bool
	Token   string
}

var instSvr *server.Server
var counter *server.Counter

func main() {

	var kubeconfig *string
	var logLevel *string
	var corsonoff *string
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stderr)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

	kubeconfig = flag.String("kubeconfig", "", "The path to the kubeconfig used to connect to the Kubernetes API server and the Kubelets (defaults to in-cluster config)")
	logLevel = flag.String("log-level", "debug", "The log level")
	corsonoff = flag.String("corsonoff", "on", "CORS(Cross-Origin Resource Sharing) on/off (defaults to on(blocked by CORS))")

	flag.Parse()

	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatal(err)
	} else {
		log.SetLevel(level)
	}

	// configuration
	//      customized by kore-board
	config.SetKubeconfig(*kubeconfig)
	config.Setup()

	//ConfigMap update delay problem quick fix..
	WatchConfig()

	//서버생성용
	instSvr, err = makeSvr()
	if err != nil {
		log.Errorf("web terminal server create error (cause=%v)", err)
	}

	// 연결 수 관리용 Counter구성
	//Todo Timeout 시간 설정
	counter = server.NewCounter(time.Duration(instSvr.Options.Timeout) * time.Second)

	r := mux.NewRouter()

	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"content-type"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
		handlers.AllowedMethods([]string{"POST, OPTIONS, GET, PUT, DELETE"}),
	)

	if *corsonoff != "on" {
		r.Use(cors)
	}

	r.HandleFunc("/api/terminal/clusters/{CLUSTER}/termtype/{TERMTYPE}", ProcTerminal).Methods("GET")
	r.HandleFunc("/api/terminal/clusters/{CLUSTER}/namespaces/{NAMESPACE}/pods/{POD}/termtype/{TERMTYPE}", ProcTerminal).Methods("GET")
	r.HandleFunc("/api/terminal/clusters/{CLUSTER}/namespaces/{NAMESPACE}/pods/{POD}/containers/{CONTAINER}/termtype/{TERMTYPE}", ProcTerminal).Methods("GET")
	r.HandleFunc("/api/terminal/ws", generateHandleWS)
	r.HandleFunc("/api/v1/config", LoadConfig).Methods("PATCH")
	r.HandleFunc("/healthy", healthy).Methods("GET") // healthy

	// Bind to a port and pass our router in
	//log.Fatal(http.ListenAndServe(":3003", handlers.CombinedLoggingHandler(os.Stdout, r)))
	log.Fatal(http.ListenAndServe(":3003", r))

	if err != nil {
		log.Errorf("web terminal server create error (cause=%v)", err)
	}
}

func LoadConfig(w http.ResponseWriter, _ *http.Request) {
	config.Setup()
	msg := fmt.Sprint("Kubeconfig updata successful")
	_, err := w.Write([]byte(msg))
	if err != nil {
		log.Errorf("Error cannot reload kubeconfig: %v", err)
	}
}

func WatchConfig() {
	v := viper.New()
	v.AutomaticEnv()
	v.SetConfigType("yaml")
	v.SetConfigFile(config.Value.ConfigLoadingRules.GetExplicitFile())

	// monitor the changes in the config file
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		log.WithField("file", e.Name).Info("Config file changed")
		config.Setup()
	})
}

func ProcTerminal(w http.ResponseWriter, r *http.Request) {

	//API요청 파라미터 파싱
	vars := mux.Vars(r)

	termreq := &termRequest{}
	termreq.cluster = NVL(vars["CLUSTER"], config.Value.DefaultContext)
	termreq.namespace = NVL(vars["NAMESPACE"], "")
	termreq.pod = NVL(vars["POD"], "")
	termreq.container = NVL(vars["CONTAINER"], "")
	termreq.termtype = NVL(vars["TERMTYPE"], "container")

	err := getContext(w, r, termreq)
	if err != nil {
		log.Errorf("%v", err)
		return
	} else {
		makeAuthToken(w, r, termreq)
	}

}

// NVL is null value logic
func NVL(str string, def string) string {
	if len(str) == 0 {
		return def
	}
	return str
}

//기존 client-go kubeconfig 정보사용
func getContext(w http.ResponseWriter, r *http.Request, req *termRequest) error {

	conf := config.Value.KubeConfig

	if config.Value.IsRunningInCluster { //In cluster mode
		req.kubeconfig = config.Value.InClusterConfig.Host
		req.kubetoken = config.Value.InClusterConfig.BearerToken
		req.inclustermode = "true"
	} else {
		var context *clientcmdapi.Context
		if conf.Contexts[req.cluster] != nil {
			context = conf.Contexts[req.cluster]
		} else {
			fmt.Println("context load fail! config.Setup() run")

			msg := fmt.Sprintf("Unable to find request Context (%s)", req)
			_, err := w.Write([]byte(msg))
			if err != nil {
				log.Errorf("Error cannot write response: %v", err)
			}
			return errors.New(msg)
		}

		termKubeConfig := clientcmdapi.NewConfig()
		termKubeConfig.Kind = conf.Kind
		termKubeConfig.APIVersion = conf.APIVersion

		termKubeConfig.Clusters[context.Cluster] = conf.Clusters[context.Cluster].DeepCopy()
		termKubeConfig.Contexts[req.cluster] = context.DeepCopy()
		termKubeConfig.CurrentContext = req.cluster
		termKubeConfig.AuthInfos[context.AuthInfo] = conf.AuthInfos[context.AuthInfo].DeepCopy()

		resultb, err := clientcmd.Write(*termKubeConfig)
		if err != nil {
			msg := fmt.Sprintf("Unable to find request Context (%s)", req)
			_, err := w.Write([]byte(msg))
			if err != nil {
				log.Errorf("Error cannot write response: %v", err)
			}
			return errors.New(msg)
		}

		req.kubeconfig = base64.StdEncoding.EncodeToString(resultb)
		req.inclustermode = "false"
	}
	return nil
}

func makeAuthToken(w http.ResponseWriter, r *http.Request, req *termRequest) {
	// 터미널 연결 식별을 위한 일회성 랜덤 문자열 생성 및 정보 설정
	token := randomstring.Generate(20)
	ttyParameter := cache.TtyParameter{}
	ttyParameter.Arg = make(map[string]string)
	setTtyValue(req, ttyParameter.Arg)

	//캐시 등록
	if err := instSvr.Cache.Add(token, &ttyParameter, cache.DefaultExpiration); err != nil {
		log.Errorf("save token and ttyParam err:%v", err)
		msg := fmt.Sprint("save token and ttyParam err")
		_, err := w.Write([]byte(msg))
		if err != nil {
			log.Errorf("Error cannot write response: %v", err)
		}
		return
	}

	// Go 데이타
	mem := Response{true, token}

	// JSON 인코딩
	jsonBytes2, err := json.Marshal(mem)
	if err != nil {
		panic(err)
	}

	// JSON 바이트를 문자열로 변경
	jsonString := string(jsonBytes2)
	fmt.Println(jsonString)

	w.Write(jsonBytes2)
}

func makeSvr() (*server.Server, error) {

	// 서버 옵션을 기본 값으로 설정
	appOptions := &server.Options{}
	if err := utils.ApplyDefaultValues(appOptions); err != nil {
		return nil, err
	}

	appOptions.PermitWrite = true

	// Backend 처리 (LocalCommand 처리용) 옵션을 기본 값으로 설정
	backendOptions := &localcommand.Options{}
	if err := utils.ApplyDefaultValues(backendOptions); err != nil {
		return nil, err
	}

	//터미널 실행시 전달되는 쉘환경 구성 스크립트
	initScript := "/opt/k3webterminal/config-namespace.sh"
	initArgs := []string{}

	// LocalCommand로 사용할 옵션과 명령어, 파라미터를 기준으로 LocalCommand Factory 생성
	factory, err := localcommand.NewFactory(initScript, initArgs, backendOptions)
	if err != nil {

		return nil, err
	}

	// Session 상태 저장을 위한 Redis 사용 여부 및 Redis 서버 연결 옵션을 기본 값으로 설정
	redisOptions := &server.RedisOptions{}
	if err := utils.ApplyDefaultValues(redisOptions); err != nil {
		return nil, err
	}

	// LocalCommand Factory와 옵션들을 기준으로 구동할 서버 인스턴스 생성
	svr, err := server.New(factory, appOptions, redisOptions)
	if err != nil {
		return nil, err
	}

	return svr, nil
}

func setTtyValue(req *termRequest, tty map[string]string) {
	e := reflect.ValueOf(req).Elem()
	fieldNum := e.NumField()
	for i := 0; i < fieldNum; i++ {
		v := e.Field(i)
		t := e.Type().Field(i)

		if v.String() != "" {
			tty[t.Name] = v.String()
		}
	}
}

func generateHandleWS(w http.ResponseWriter, r *http.Request) {
	fmt.Println("websocket connect")
	c := r.Context()

	ctx, cancel := context.WithCancel(c)

	go func() {
		select {
		case <-counter.Timer().C:
			cancel()
		case <-ctx.Done():
		}
	}()

	num := counter.Add(1)
	closeReason := "unknown reason"

	defer func() {
		num := counter.Done()

		log.Info(fmt.Sprintf("Connection closed: %s, %s, connections: %d/%d", r.RemoteAddr, closeReason, num, instSvr.Options.MaxConnection))
		if instSvr.Options.Once {
			cancel()
		}
	}()

	log.Info(fmt.Sprintf("New client connected: %s, connections: %d/%d", r.RemoteAddr, num, instSvr.Options.MaxConnection))

	instSvr.Upgrader.ReadBufferSize = webtty.MaxBufferSize
	instSvr.Upgrader.WriteBufferSize = webtty.MaxBufferSize
	instSvr.Upgrader.EnableCompression = true
	conn, err := instSvr.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		closeReason = err.Error()
		return
	}
	defer conn.Close()
	conn.SetCompressionLevel(9)
	err = processWSConn(ctx, conn)

	switch err {
	case ctx.Err():
		closeReason = "cancelation"
	case webtty.ErrSlaveClosed:
		closeReason = instSvr.Factory.Name()
	case webtty.ErrMasterClosed:
		closeReason = "client close"
	case webtty.ErrConnectionLostPing:
		closeReason = webtty.ErrConnectionLostPing.Error()
	default:
		closeReason = fmt.Sprintf("an error: %s", err)
	}

}

// processWSConn - 터미널과 연결할 WebSocket 연결 구성
func processWSConn(ctx context.Context, conn *websocket.Conn) error {
	fmt.Println("processWSConn")
	typ, initLine, err := conn.ReadMessage()
	if err != nil {
		return errors.Wrapf(err, "failed to authenticate websocket connection")
	}
	if typ != websocket.TextMessage {
		return errors.New("failed to authenticate websocket connection: invalid message type")
	}

	//클라이언트에서 전달받은 토큰처리
	var init server.InitMessage
	err = json.Unmarshal(initLine, &init)
	if err != nil {
		return errors.Wrapf(err, "failed to authenticate websocket connection")
	}

	params := map[string]string{}

	if len(init.AuthToken) > 0 {
		ttyParameter := instSvr.Cache.Get(init.AuthToken)
		cachedKey := init.AuthToken

		if ttyParameter != nil {
			params = ttyParameter.Arg
			instSvr.Cache.Delete(cachedKey)
		} else {
			return errors.New("ERROR:Invalid Token")
		}
	} else {
		return errors.New("ERROR:No Token Provided")
	}

	//Backend Slave생성
	var slave server.Slave
	slave, err = instSvr.Factory.New(params)
	if err != nil {
		return errors.Wrapf(err, "failed to create backend")
	}
	defer slave.Close()

	//Webtty 설정
	opts := []webtty.Option{}

	if instSvr.Options.PermitWrite {
		opts = append(opts, webtty.WithPermitWrite())
	}
	if instSvr.Options.EnableReconnect {
		opts = append(opts, webtty.WithReconnect(instSvr.Options.ReconnectTime))
	}
	if instSvr.Options.Width > 0 {
		opts = append(opts, webtty.WithFixedColumns(instSvr.Options.Width))
	}
	if instSvr.Options.Height > 0 {
		opts = append(opts, webtty.WithFixedRows(instSvr.Options.Height))
	}
	if instSvr.Options.Preferences != nil {
		opts = append(opts, webtty.WithMasterPreferences(instSvr.Options.Preferences))
	}

	tty, err := webtty.New(&server.WsWrapper{conn}, slave, opts...)
	if err != nil {
		return errors.Wrapf(err, "failed to create webtty")
	}

	err = tty.Run(ctx)

	return err
}

func healthy(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("healthy")
}
