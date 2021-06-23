package config

import (
	"flag"
	"os"
	"strings"

	"github.com/acornsoftlab/dashboard/pkg/auth"
	"github.com/acornsoftlab/dashboard/pkg/lang"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
)

var Value = &conf{}
var Authenticator *auth.Authenticator // authentication/authroization config
var Cluster *kubeCluster

func init() {

	// startup parameters
	logLevel := flag.String("log-level", os.Getenv("LOG_LEVEL"), "The log level")
	flag.StringVar(&Value.MetricsScraperUrl, "metrics-scraper-url", os.Getenv("METRICS_SCRAPER_URL"), "The address of the metrics-scraper rest-api URL")
	flag.StringVar(&Value.TerminalUrl, "terminal-url", os.Getenv("TERMINAL_URL"), "The address of the Terminal server")
	kubeconfig := flag.String("kubeconfig", "", "The path to the kubeconfig used to connect to the Kubernetes API server and the Kubelets")
	authconfig := flag.String("auth", os.Getenv("AUTH"), "The authenticate options")

	//k8s.io client-go logs
	flag.Set("logtostderr", "ture")
	flag.Set("stderrthreshold", "FATAL")

	flag.Parse()

	//set default
	//*kubeconfig = "strategy=configmap,configmap=kore-board-kubeconfig,namespace=kore,filename=config"
	*authconfig = lang.NVL(*authconfig, "strategy=cookie,secret=static-token,token=acornsoft")
	Value.MetricsScraperUrl = lang.NVL(Value.MetricsScraperUrl, "http://localhost:8000")
	Value.TerminalUrl = lang.NVL(Value.TerminalUrl, "http://localhost:3003")
	*logLevel = lang.NVL(*logLevel, "debug")

	//logger
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stderr)

	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatal(err)
	} else {
		log.SetLevel(level)
		log.Infof("Log level is '%s'", *logLevel)
	}

	// print startup parameters
	log.Infof("Startup parameter 'metrics-scraper-url' is '%s'", Value.MetricsScraperUrl)
	log.Infof("Startup parameter 'kubeconfig' is '%s'", *kubeconfig)
	log.Infof("Startup parameter 'auth' is '%s'", *authconfig)

	// unmarshall kubeconfig
	Value.KubeConfig = &kubeConfig{}
	Value.KubeConfig.Data = make(map[string]string)
	if *kubeconfig == "" || !strings.Contains(*kubeconfig, "strategy=") {
		Value.KubeConfig.Strategy = "file"
		Value.KubeConfig.Data["path"] = *kubeconfig
	} else {
		for _, e := range strings.Split(*kubeconfig, ",") {
			parts := strings.Split(e, "=")
			if parts[0] == "strategy" {
				Value.KubeConfig.Strategy = parts[1]
			} else {
				Value.KubeConfig.Data[parts[0]] = parts[1]
			}
		}
	}

	// unmarshall auth-config
	Value.AuthConfig = &auth.AuthConfig{}
	Value.AuthConfig.Data = make(map[string]string)
	for _, e := range strings.Split(*authconfig, ",") {
		parts := strings.Split(e, "=")
		if parts[0] == "strategy" {
			Value.AuthConfig.Strategy = parts[1]
		} else if parts[0] == "secret" {
			Value.AuthConfig.Secret = parts[1]
		} else if parts[0] == "access-key" {
			Value.AuthConfig.AccessKey = parts[1]
		} else if parts[0] == "refresh-key" {
			Value.AuthConfig.RefreshKey = parts[1]
		} else {
			Value.AuthConfig.Data[parts[0]] = parts[1]
		}
	}

}

// 재로딩 가능한  config 정의
func Setup() {
	var err error

	// kubeconfig 파일 로드
	if Cluster, err = newKubeCluster(*Value.KubeConfig); err != nil {
		log.Errorf("can't create kubernetes clusters (cause=%s)", err.Error())
	} else {
		Cluster.IsRunningInCluster = (len(Cluster.ClusterNames) == 1 && Cluster.InCluster != nil)
		if len(Cluster.ClusterNames) == 0 {
			log.Warnf("Initialized empty clusters (kubeconfig=none, in-cluster=none, running-in-cluster=none)")
		} else {
			log.Infof("Initialized clusters (kubeconfig-strategy=%s, in-cluster=%t, running-in-cluster=%t, contexts=%s, default-context=%s)", Value.KubeConfig.Strategy, (Cluster.InCluster != nil), Cluster.IsRunningInCluster, Cluster.ClusterNames, Cluster.DefaultContext)
		}
	}

	// definition "authenticator"
	if Cluster.DefaultContext != "" && Value.AuthConfig != nil {
		var restConfig *rest.Config
		if Cluster.InCluster != nil {
			restConfig = Cluster.InCluster.RESTConfig
		}
		if Authenticator, err = auth.CreateAuthenticator(Value.AuthConfig, restConfig); err != nil {
			log.Errorf("Invalid authenticator (cause=%s)", err.Error())
		}
	}
	if Authenticator == nil {
		Authenticator = auth.DummyAuthenticator()
		log.Warnln("Initialized dummy authenticator")
	} else {
		log.Infof("Initialized authenticator (strategy=%s, provider=%s)", Value.AuthConfig.Strategy, Value.AuthConfig.Secret)
	}

}
