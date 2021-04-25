package config

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/acornsoftlab/dashboard/pkg/lang"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"os"

	"path/filepath"
)

const (
	IN_CLUSTER_NAME = "kubernetes@in-cluster"
)

type conf struct {
	ConfigLoadingRules clientcmd.ClientConfigLoader // kubeconfig loading rule
	KubeConfig         *api.Config                  // kubeconfig file
	DefaultContext     string                       // kubeconfig file - default context
	Contexts           []string                     // kubeconfig file - context list
	InClusterConfig    *rest.Config                 // api clietner kubeconfig - in-cluter
	kubeConfigs        map[string]*rest.Config      // api clietner kubeconfig - each context
	MetricsScraperUrl  *string                      // metrics scraper
	SecretToken        string                       // secret token
	SecretAccessKey    *string                      // access-key singing secret
	SecretRefreshKey   *string                      // refresh-key singing secret
	IsRunningInCluster bool                         // running in-cluster mode
}

var Value = &conf{}

var kubeconfig *string
var secretTokenPath *string
var loader clientcmd.ClientConfigLoader

func NewKubeClient(ctx string) (*kubernetes.Clientset, error) {
	conf, err := KubeConfigs(ctx)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(conf)

}

func KubeConfigs(ctx string) (*rest.Config, error) {
	if lang.ArrayContains(Value.Contexts, ctx) {
		return Value.kubeConfigs[ctx], nil
	} else {
		return nil, errors.New(fmt.Sprintf("cannot found context=%s", ctx))
	}
}
func KubeConfigsDefault(ctx string) *rest.Config {
	if lang.ArrayContains(Value.Contexts, ctx) {
		return Value.kubeConfigs[ctx]
	} else {
		return Value.kubeConfigs[Value.DefaultContext]
	}
}

func init() {

	// arguments
	kubeconfig = flag.String("kubeconfig", "", "The path to the kubeconfig used to connect to the Kubernetes API server and the Kubelets (defaults to in-cluster config)")
	Value.MetricsScraperUrl = flag.String("metrics-scraper-url", "http://localhost:8000", "The address of the Metrics-scraper Apiserver")
	logLevel := flag.String("log-level", "debug", "The log level")
	secretTokenPath = flag.String("token", os.Getenv("TOKEN"), "The secret token path")
	Value.SecretAccessKey = flag.String("access-secret", "whdmstkddkzhsthvmxm", "The singing secret of access key")
	Value.SecretRefreshKey = flag.String("refresh-ecret", "dkzhsthvmxmwhdmstkd", "The singing secret of refresh key")

	//logger
	err := flag.Set("logtostderr", "true")
	if err != nil {
		log.Errorf("Error cannot set logtostderr: %v", err)
	}
	flag.Parse()

	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stderr)

	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatal(err)
	} else {
		log.SetLevel(level)
	}

	// KUBECONFIG 로드
	//  1순위. "--kubeconfig" 옵션에서 로드한다.
	//  2순위. env "KUBECONFIG" 값으로 로드한다.
	//  3순위."~/.kube/config" 에서 로드한다.
	//  4순위. in-cluster-config 로드한다.
	if *kubeconfig != "" { // load from --kubeconfig
		Value.ConfigLoadingRules = &clientcmd.ClientConfigLoadingRules{ExplicitPath: *kubeconfig}
	} else {
		Value.ConfigLoadingRules = clientcmd.NewDefaultClientConfigLoadingRules()
	}
}

// 재로딩 가능한  config 정의
func Setup() {
	var err error

	// kubeconfig 파일 로드
	Value.kubeConfigs = map[string]*rest.Config{}
	Value.Contexts = []string{}

	// in-cluster
	Value.InClusterConfig, err = rest.InClusterConfig()
	if err != nil {
		log.Warnln("cannot load kubeconfig in-cluster")
	}

	// kubeconfig 파일 로드
	Value.KubeConfig, err = Value.ConfigLoadingRules.Load()
	if err != nil {
		log.Warnf("cannot load kubeconfig: %s (cause=%v)", *kubeconfig, err)
	} else {
		for key := range Value.KubeConfig.Contexts {
			contextCfg, err := clientcmd.NewNonInteractiveClientConfig(*Value.KubeConfig, key, &clientcmd.ConfigOverrides{}, loader).ClientConfig()
			if err == nil {
				Value.Contexts = append(Value.Contexts, key)
				Value.kubeConfigs[key] = contextCfg
			}
		}
		Value.DefaultContext = Value.KubeConfig.CurrentContext
	}

	// 로드된 context가 없다면 in-cluster
	Value.IsRunningInCluster = (len(Value.Contexts) == 0 && Value.InClusterConfig != nil)
	if Value.IsRunningInCluster {
		Value.kubeConfigs[IN_CLUSTER_NAME] = Value.InClusterConfig
		Value.Contexts = []string{IN_CLUSTER_NAME}
	}

	// print infomation
	if Value.InClusterConfig == nil {
		log.Infof("contexts (len=%v) , non-in-cluster", len(Value.Contexts))
	} else {
		log.Infof("contexts (len=%v) , in-cluster OK", len(Value.Contexts))
	}

	// default-context = 없다면 첫번째
	if len(Value.Contexts) > 0 && Value.DefaultContext == "" {
		Value.DefaultContext = Value.Contexts[0]
	}

	// authenticate 를 위한 secret token 초기화
	initScretToken(*secretTokenPath)

	// print infomation
	log.Infof("TOKEN : %s", Value.SecretToken)

}

func initScretToken(path string) {

	var apiClient *kubernetes.Clientset
	var err error

	// secret 을 읽어올 cluster 선정
	if Value.InClusterConfig == nil {
		apiClient, err = NewKubeClient(Value.DefaultContext)
		if err != nil {
			log.Errorln(err.Error())
		}
	} else {
		apiClient, err = kubernetes.NewForConfig(Value.InClusterConfig)
		if err != nil {
			log.Errorln(err.Error())
		}
	}

	if apiClient != nil {
		ns := lang.NVL(os.Getenv("NAMESPACE"), "kore")
		nm := lang.NVL(os.Getenv("SERVICE_ACCOUNT"), "kore-board")
		sa, err := apiClient.CoreV1().ServiceAccounts(ns).Get(context.TODO(), nm, metav1.GetOptions{})
		if err != nil {
			log.Warnf("cannot load service account env.NAMESPACE=%s, env.SERVICE_ACCOUNT=%s, (cause %s)", ns, nm, err)
		} else {
			se, err := apiClient.CoreV1().Secrets(ns).Get(context.TODO(), sa.Secrets[0].Name, metav1.GetOptions{})
			if err == nil {
				Value.SecretToken = string(se.Data["token"])
			} else {
				log.Warnf("cannot load token env.NAMESPACE=%s, env.SERVICE_ACCOUNT=%s, (cause %s)", ns, nm, err)
			}
		}
	} else {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Infof("can not find a secret token file '%s'", path)
		} else {
			if d, err := ioutil.ReadFile(path); err == nil {
				Value.SecretToken = string(d)
				log.Infof("read a secret token from '%s'", path)
			} else {
				log.Errorf("cannot read a secret token from '%s' (cause=%s)", path, err)
			}
		}
	}

	if Value.SecretToken == "" {
		b := make([]byte, 256) //490
		rand.Read(b)
		Value.SecretToken = fmt.Sprintf("%x", b)
		log.Warnln("generate a new random secret token")
	}

	// save token file
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, os.FileMode(777)); err != nil {
			log.Errorf("can not make a directory (dir=%s, cause=%s)", dir, err.Error())
		} else {
			log.Infof("create a token directory (dir=%s)", dir)
		}
	}

	os.Remove(path)
	b := []byte(Value.SecretToken)
	if err := ioutil.WriteFile(path, b, os.FileMode(777)); err != nil {
		log.Errorf("cannot write a secret token file (cause=%s, path=%s)", err.Error(), path)
	} else {
		log.Infof("create a token file (path=%s)", path)
	}

}
