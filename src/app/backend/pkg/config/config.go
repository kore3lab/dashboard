package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	"github.com/acornsoftlab/dashboard/pkg/auth"
	"github.com/acornsoftlab/dashboard/pkg/lang"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

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

func init() {

	// arguments
	kubeconfig = flag.String("kubeconfig", "", "The path to the kubeconfig used to connect to the Kubernetes API server and the Kubelets (defaults to in-cluster config)")
	flag.StringVar(&Value.MetricsScraperUrl, "metrics-scraper-url", "http://localhost:8000", "The address of the Metrics-scraper Apiserver")
	logLevel := flag.String("log-level", "debug", "The log level")
	authconfig := flag.String("auth", auth.DefaultAuthenticator, fmt.Sprintf("The authenticator (default=%s)", auth.DefaultAuthenticator))

	//logger
	err := flag.Set("logtostderr", "true")
	if err != nil {
		log.Warnf("error cannot set logtostderr: %s", err.Error())
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

	// definition "authenticator"
	log.Infof("auth-config %v", *authconfig)
	if err = json.Unmarshal([]byte(*authconfig), &Value.AuthConfig); err != nil {
		log.Errorf("invalid auth parameter (cause=%s)", err.Error())
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
		log.Infof("contexts %v , non-in-cluster", Value.Contexts)
	} else {
		log.Infof("contexts %v , in-cluster OK", Value.Contexts)
	}

	// default context (1:in-cluster > 2: first context)
	if Value.IsRunningInCluster {
		Value.DefaultContext = IN_CLUSTER_NAME
	} else if len(Value.Contexts) > 0 {
		Value.DefaultContext = Value.Contexts[0]
	}
	if Value.DefaultContext == "" {
		log.Warnln("there is no default context")
	} else {
		log.Infof("default context is \"%s\"", Value.DefaultContext)
	}

	// definition "authenticator"
	if Value.AuthConfig != nil {
		if Value.Authenticator, err = auth.CreateAuthenticator(Value.AuthConfig, Value.kubeConfigs[Value.DefaultContext]); err != nil {
			log.Errorf("invalid authenticator (cause=%s)", err.Error())
		}
	}
	if Value.Authenticator == nil {
		Value.Authenticator = auth.DummyAuthenticator()
		log.Warnln("initialized dummy authenticator")
	} else {
		log.Infof("initialized authenticator (strategy=%s, provider=%s)", Value.AuthConfig.Strategy, Value.AuthConfig.Secret["type"])
	}

}
