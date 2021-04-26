package config // add by kore-board

import (
	"errors"
	"fmt"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	log "github.com/sirupsen/logrus"
)

const (
	IN_CLUSTER_NAME = "kubernetes@in-cluster"
)

type conf struct {
	ConfigLoadingRules clientcmd.ClientConfigLoader
	DefaultContext     string
	kubeConfigs        map[string]*rest.Config
	Contexts           []string
}

var Value = &conf{}

func KubeConfigs(ctx string) (*rest.Config, error) {

	for _, s := range Value.Contexts {
		if s == ctx {
			return Value.kubeConfigs[ctx], nil
		}
	}

	return nil, errors.New(fmt.Sprintf("cannot found context=%s", ctx))
}

func SetKubeconfig(kubeconfig string) {
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
	Value.ConfigLoadingRules = loader
}

func Setup() {

	Value.kubeConfigs = map[string]*rest.Config{}
	Value.Contexts = []string{}

	// kubeconfig 파일 로드
	cfg, err := Value.ConfigLoadingRules.Load()
	if err == nil {
		for key := range cfg.Contexts {
			contextCfg, err := clientcmd.NewNonInteractiveClientConfig(*cfg, key, &clientcmd.ConfigOverrides{}, Value.ConfigLoadingRules).ClientConfig()
			if err == nil {
				Value.Contexts = append(Value.Contexts, key)
				Value.kubeConfigs[key] = contextCfg
			}
		}
		Value.DefaultContext = cfg.CurrentContext
	} else {
		log.Warnf("cannot load kubeconfig (cause=%v)", err)
	}

	// 로드된 context가 없다면 in-cluster 모드
	if len(Value.Contexts) == 0 {
		cnf, err := rest.InClusterConfig()
		if err != nil {
			log.Error("running empty cluster (cannot load a kubeconfig in-cluster)")
		} else {
			Value.kubeConfigs[IN_CLUSTER_NAME] = cnf
			Value.Contexts = []string{IN_CLUSTER_NAME}
			log.Infoln("running in-cluster mode")
		}
	} else {
		log.Infof("running out-cluster mode (contexts=%v)", len(Value.Contexts))
	}

	if len(Value.Contexts) > 0 && Value.DefaultContext == "" {
		Value.DefaultContext = Value.Contexts[0]
	}

}
