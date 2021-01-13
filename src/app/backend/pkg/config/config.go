package config

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type conf struct {
	CurrentContext string
	KubeConfigs    map[string]*rest.Config
	Contexts       []string
}

var Value = &conf{}

func Setup(kubeconfig string) {
	// KUBECONFIG 로드 룰
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

	Value.KubeConfigs = map[string]*rest.Config{}

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

	}

	// 로드된 context가 없다면 in-cluster 모드
	if len(Value.Contexts) == 0 {
		cnf, err := rest.InClusterConfig()
		if err != nil {
			panic("cannot load kubeconfig inCluster")
		}
		Value.KubeConfigs["default"] = cnf
		Value.Contexts = []string{"default"}
	}

	if Value.CurrentContext == "" {
		Value.CurrentContext = Value.Contexts[0]
	}

}
