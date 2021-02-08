package config // add by acornsoft-dashboard

import (
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type conf struct {
	Kubeconfig string
	CurrentContext string
	KubeConfigs    map[string]*rest.Config
	Contexts       []string
}

var Value = &conf{}

func Setup(kubeconfig string) {

	// KUBECONFIG 로드 Rules
	//  1순위. "--kubeconfig" 옵션에서 로드한다.
	//  2순위. env "KUBECONFIG" 값으로 로드한다.
	//  3순위."~/.kube/config" 에서 로드한다.
	//  4순위. in-cluster-config 로드한다.

	var loader clientcmd.ClientConfigLoader
	if kubeconfig != "" { // load from --kubeconfig
		loader = &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig}
	} else {
		loader = clientcmd.NewDefaultClientConfigLoadingRules()
	}

	Value.Kubeconfig = kubeconfig
	Value.KubeConfigs = map[string]*rest.Config{}
	Value.Contexts = []string{}

	// kubeconfig 파일 로드
	cfg, err := loader.Load()
	if err == nil {
		for key := range cfg.Contexts {
			contextCfg, err := clientcmd.NewNonInteractiveClientConfig(*cfg, key, &clientcmd.ConfigOverrides{}, loader).ClientConfig()
			if err == nil {
				Value.Contexts = append(Value.Contexts, key)
				Value.KubeConfigs[key] = contextCfg
			}
		}
		Value.CurrentContext = cfg.CurrentContext
	} else {
		log.Warnf("cannot load kubeconfig: %s (cause=%v)", kubeconfig, err)
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
