package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

type conf struct {
	ConfigLoadingRules clientcmd.ClientConfigLoader
	KubeConfig         *api.Config
	CurrentContext     string
	KubeConfigs        map[string]*rest.Config
	Contexts           []string
	MetricsScraperUrl  string
}

var Value = &conf{}

var kubeconfig *string
var loader clientcmd.ClientConfigLoader

func init() {
	var metricsScraperUrl *string
	var logLevel *string

	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)

	// arguments
	kubeconfig = flag.String("kubeconfig", "", "The path to the kubeconfig used to connect to the Kubernetes API server and the Kubelets (defaults to in-cluster config)")
	metricsScraperUrl = flag.String("metrics-scraper-url", "http://localhost:8000", "The address of the Metrics-scraper Apiserver")
	logLevel = flag.String("log-level", "info", "The log level")

	//logger
	err := flag.Set("logtostderr", "true")
	if err != nil {
		log.Errorf("Error cannot set logtostderr: %v", err)
	}
	flag.Parse()

	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatal(err)
	} else {
		log.SetLevel(level)
	}

	//metrics-scraper-host
	Value.MetricsScraperUrl = *metricsScraperUrl

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

	// kubeconfig 파일 로드
	var err error
	Value.KubeConfig, err = Value.ConfigLoadingRules.Load()
	if err != nil {
		log.Warnf("cannot load kubeconfig: %s (cause=%v)", kubeconfig, err)
	}

	Value.KubeConfigs = map[string]*rest.Config{}
	Value.Contexts = []string{}

	// kubeconfig 파일 로드
	if Value.KubeConfig != nil {
		for key := range Value.KubeConfig.Contexts {
			contextCfg, err := clientcmd.NewNonInteractiveClientConfig(*Value.KubeConfig, key, &clientcmd.ConfigOverrides{}, loader).ClientConfig()
			if err == nil {
				Value.Contexts = append(Value.Contexts, key)
				Value.KubeConfigs[key] = contextCfg
			}
		}
		Value.CurrentContext = Value.KubeConfig.CurrentContext
	}

	// 로드된 context가 없다면 in-cluster 모드
	if len(Value.Contexts) == 0 {
		cnf, err := rest.InClusterConfig()
		if err != nil {
			log.Panic("cannot load kubeconfig inCluster")
		}
		Value.KubeConfigs["default"] = cnf
		Value.Contexts = []string{"default"}
	}

	if Value.CurrentContext == "" {
		Value.CurrentContext = Value.Contexts[0]
	}

}
