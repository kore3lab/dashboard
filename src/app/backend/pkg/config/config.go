package config

import (
	"os"

	"github.com/acornsoftlab/kore3/pkg/lang"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	// clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type conf struct {
	KubeConfigPath string
	CurrentContext string
	KubeConfigs    map[string]*rest.Config
	Contexts       []string
}

var Value = &conf{}

func Setup() {
	// "GIN_MODE"
	//      "debug", "release", "test"
	//      https://github.com/gin-gonic/gin/blob/master/mode.go
	Value.KubeConfigPath = lang.NVL(os.Getenv("KUBECONFIG"), os.Getenv("HOME")+"/.kube/config")

	// kubeconf Context
	kubeconfig := Value.KubeConfigPath

	var loader clientcmd.ClientConfigLoader
	if kubeconfig != "" { // load from --kubeconfig
		loader = &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig}
	} else {
		loader = clientcmd.NewDefaultClientConfigLoadingRules()
	}

	cfg, err := loader.Load()
	if err != nil {
		panic("cannot load kubecfg")
	}

	configs := map[string]*rest.Config{}
	for context := range cfg.Contexts {
		contextCfg, err := clientcmd.NewNonInteractiveClientConfig(*cfg, context, &clientcmd.ConfigOverrides{}, loader).ClientConfig()
		if err != nil {
			panic("could not create client config")
		}
		configs[context] = contextCfg
		Value.Contexts = append(Value.Contexts, context)
	}

	Value.KubeConfigs = configs
	Value.CurrentContext = cfg.CurrentContext

}
