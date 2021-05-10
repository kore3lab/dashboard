package config

import (
	"github.com/acornsoftlab/dashboard/pkg/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
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
	MetricsScraperUrl  string                       // metrics scraper
	IsRunningInCluster bool                         // running in-cluster mode
	AuthConfig         *auth.AuthConfig             // auth-config
	Authenticator      *auth.Authenticator          // authentication/authroization config

	// SecretToken        string                       // secret token
	//AccessTokenSecret string // access-key singing secret
	//RefreshToekSecret string // refresh-key singing secret
}
