package config

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/acornsoftlab/dashboard/pkg/client"
	"github.com/acornsoftlab/dashboard/pkg/lang"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

func (me *kubeCluster) Create(apiConfig api.Config) error {

	me.KubeConfig = apiConfig.DeepCopy()
	return me.Save()

}
func (me *kubeCluster) Add(name string, json map[string]interface{}) error {

	// create objects
	cluster := &api.Cluster{}
	context := &api.Context{}
	user := &api.AuthInfo{}

	// context, cluster, user 이름 중복 회피
	if me.KubeConfig != nil {
		if _, exist := me.KubeConfig.Contexts[name]; exist {
			name = fmt.Sprintf("%s.%s", name, lang.RandomString(3))
		}
		context.Cluster = fmt.Sprintf("%s-cluster", name)
		context.AuthInfo = fmt.Sprintf("%s-user", name)

		if _, exist := me.KubeConfig.Clusters[context.Cluster]; exist {
			context.Cluster = fmt.Sprintf("%s-%s-cluster", name, lang.RandomString(5))
		}
		if _, exist := me.KubeConfig.AuthInfos[context.AuthInfo]; exist {
			context.AuthInfo = fmt.Sprintf("%s-%s-user", name, lang.RandomString(5))
		}
	} else {
		context.Cluster = fmt.Sprintf("%s-cluster", name)
		context.AuthInfo = fmt.Sprintf("%s-user", name)
	}

	// parsing
	jsonC := json["cluster"].(map[string]interface{})
	jsonU := json["user"].(map[string]interface{})

	// parsing - cluster
	cluster.Server = jsonC["server"].(string)
	val, exists := jsonC["certificate-authority-data"].(string)
	if exists {
		ca, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			return fmt.Errorf("Unable to decode cerificate-authority-data (data=%s, cause=%s)", val, err)
		}
		cluster.CertificateAuthorityData = ca
	} else {
		cluster.CertificateAuthority = jsonC["certificate-authority"].(string)
	}

	// parsing - user
	val, exists = jsonU["client-certificate-data"].(string)

	if exists {
		ca, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			return fmt.Errorf("Unable to decode client-certificate-data (data=%s, cuase=%s)", val, err)
		}
		user.ClientCertificateData = ca
	} else {
		user.ClientCertificate = jsonU["client-certificate"].(string)
	}

	val, exists = jsonU["client-key-data"].(string)
	if exists {
		ca, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			return fmt.Errorf("Unable to decode client-key-data (data=%s, cause=%s)", val, err)
		}
		user.ClientKeyData = ca
	} else {
		user.ClientKey = jsonU["client-key"].(string)
	}

	me.KubeConfig.Clusters[context.Cluster] = cluster
	me.KubeConfig.AuthInfos[context.AuthInfo] = user
	me.KubeConfig.Contexts[name] = context

	return me.Save()

}

func (me *kubeCluster) Remove(name string) error {

	conf := me.KubeConfig.DeepCopy()

	if conf.Contexts[name] != nil {
		if conf.Clusters[conf.Contexts[name].Cluster] != nil {
			delete(conf.Clusters, conf.Contexts[name].Cluster)
		}
		if conf.AuthInfos[conf.Contexts[name].AuthInfo] != nil {
			delete(conf.AuthInfos, conf.Contexts[name].AuthInfo)
		}
		if conf.CurrentContext == name {
			conf.CurrentContext = ""
		}
		delete(conf.Contexts, name)

	} else {
		return fmt.Errorf("not found context %s", name)
	}

	me.KubeConfig = conf
	return me.Save()

}

func newKubeCluster(conf kubeConfig) (*kubeCluster, error) {

	// cluster
	var provider *cluserConfigProvider
	var err error

	if conf.Strategy == KubeConfigStrategyConfigmap {
		provider, err = createConfigmapClusterConfig(conf.Data["namespace"], conf.Data["configmap"], conf.Data["filename"])
		if err != nil {
			log.Warnf("can't create a 'configmap' kubeconfig-provider (cause=%s)", err.Error())
		}
	}

	// default provider
	if provider == nil {
		if provider, err = createFileClusterConfig(conf.Data["path"]); err != nil {
			log.Warnf("can't create a 'file' kubeconfig-provider  (cause=%s)", err.Error())
		}
	}
	if provider == nil {
		log.Panicf("can't create a kubeconfig-provider")
	}

	// properties
	clusters := make(map[string]*kubeClient)
	clusterNames := []string{}
	defaultContext := ""

	apiConfig, err := provider.load()
	if err != nil {
		log.Warnf("can't load a kubeconfig (cause=%s)", err.Error())
		apiConfig = &api.Config{
			Clusters:  make(map[string]*api.Cluster),
			AuthInfos: make(map[string]*api.AuthInfo),
			Contexts:  make(map[string]*api.Context),
		}
	} else {

		// "Context", "KubeConfigs"
		if apiConfig.CurrentContext != "" {
			defaultContext = apiConfig.CurrentContext
		}
		for key := range apiConfig.Contexts {
			if defaultContext == "" {
				defaultContext = key
			}

			if restConfig, err := clientcmd.NewNonInteractiveClientConfig(*apiConfig, key, &clientcmd.ConfigOverrides{}, nil).ClientConfig(); err == nil {
				clusterNames = append(clusterNames, key)
				clusters[key] = createKubeClient(key, restConfig)
			}
		}
	}

	// in-cluster
	inClusterConfig, _ := rest.InClusterConfig()

	// 로드된 context가 없다면 in-cluster 사용
	if len(clusters) == 0 && inClusterConfig != nil {
		clusters[IN_CLUSTER_NAME] = createKubeClient(IN_CLUSTER_NAME, inClusterConfig)
		clusterNames = []string{IN_CLUSTER_NAME}
		defaultContext = IN_CLUSTER_NAME
	}

	// kubeCluster
	c := &kubeCluster{
		KubeConfig:     apiConfig,
		clients:        clusters,
		ClusterNames:   clusterNames,
		DefaultContext: defaultContext,
	}
	if inClusterConfig != nil {
		c.InCluster = createKubeClient(IN_CLUSTER_NAME, inClusterConfig)
	}
	c.IsRunningInCluster = (len(c.ClusterNames) == 1 && c.InCluster != nil) // running within a cluster

	// kubeCluster.Client()
	c.Client = func(context string) (*kubeClient, error) {
		if !lang.ArrayContains(c.ClusterNames, context) {
			if context == IN_CLUSTER_NAME && c.InCluster != nil {
				return c.InCluster, nil
			} else {
				return nil, fmt.Errorf("can't find a context '%s' in %v", context, c.ClusterNames)
			}
		}
		return c.clients[context], nil
	}

	// kubeCluster.Save()
	c.Save = func() error {
		return provider.save(c.KubeConfig)
	}

	return c, nil
}

type cluserConfigProvider struct {
	load func() (*api.Config, error)
	save func(*api.Config) error
}

type fileCluserConfigProvider struct {
	cluserConfigProvider
	defaultFilename string
}

// kubeconfig - file
func createFileClusterConfig(kubeconfig string) (*cluserConfigProvider, error) {

	conf := &fileCluserConfigProvider{}

	conf.load = func() (*api.Config, error) {
		var configLoadingRules clientcmd.ClientConfigLoader
		if kubeconfig != "" {
			configLoadingRules = &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig}
		} else {
			configLoadingRules = clientcmd.NewDefaultClientConfigLoadingRules()
		}
		apiConfig, err := configLoadingRules.Load()
		conf.defaultFilename = configLoadingRules.GetDefaultFilename()
		if err != nil {
			return nil, err
		}

		return apiConfig.DeepCopy(), nil

	}
	conf.save = func(apiConfig *api.Config) error {
		if err := clientcmd.WriteToFile(*apiConfig, conf.defaultFilename); err != nil {
			return err
		}
		return nil
	}

	return &conf.cluserConfigProvider, nil

}

// kubeconfig - configmap
func createConfigmapClusterConfig(nsName string, cmName string, fileName string) (*cluserConfigProvider, error) {
	conf := &cluserConfigProvider{}

	conf.load = func() (*api.Config, error) {
		rest, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}

		clientset, err := kubernetes.NewForConfig(rest)
		cm, err := clientset.CoreV1().ConfigMaps(nsName).Get(context.TODO(), cmName, v1.GetOptions{})
		if err != nil {
			return nil, err
		}

		if cm.Data[fileName] == "" {
			return nil, fmt.Errorf("kubeconfig data is empty namespace=%s, configmap=%s, filename=%s, data=%v", nsName, cmName, fileName, cm.Data)
		}
		clientConfig, err := clientcmd.NewClientConfigFromBytes([]byte(cm.Data[fileName]))
		if err != nil {
			return nil, err
		}

		//apiConfig := &api.Config{}
		apiConfig, err := clientConfig.RawConfig()
		if err != nil {
			return nil, err
		}

		return apiConfig.DeepCopy(), nil
	}

	conf.save = func(apiConfig *api.Config) error {
		rest, err := rest.InClusterConfig()
		if err != nil {
			return err
		}

		clientset, err := kubernetes.NewForConfig(rest)
		cm, err := clientset.CoreV1().ConfigMaps(nsName).Get(context.TODO(), cmName, v1.GetOptions{})
		if err != nil {
			return err
		}

		b, err := clientcmd.Write(*apiConfig)
		if err != nil {
			return err
		}
		cm.Data[fileName] = string(b)

		_, err = clientset.CoreV1().ConfigMaps(nsName).Update(context.TODO(), cm, v1.UpdateOptions{})
		if err != nil {
			return err
		}
		return nil
	}

	return conf, nil
}

func createKubeClient(name string, restConfig *rest.Config) *kubeClient {

	cluster := &kubeClient{Name: name, RESTConfig: restConfig}

	// NewDynamicClient
	cluster.NewMetricsClient = func() (*versioned.Clientset, error) {
		return versioned.NewForConfig(restConfig)
	}
	// NewKubernetesClient
	cluster.NewKubernetesClient = func() (*kubernetes.Clientset, error) {
		return kubernetes.NewForConfig(restConfig)
	}

	// NewDiscoveryClient
	cluster.NewDiscoveryClient = func() (*discovery.DiscoveryClient, error) {
		return discovery.NewDiscoveryClientForConfig(restConfig)
	}

	// NewDynamicClient
	cluster.NewDynamicClient = func() (*client.DynamicClient, error) {
		return client.NewDynamicClient(restConfig), nil
	}

	// 예:  schema.GroupVersionResource{Group: "networking.istio.io", Version: "v1alpha3", Resource: "virtualservices"}
	cluster.NewDynamicClientSchema = func(group string, version string, resource string) (*client.DynamicClient, error) {
		return client.NewDynamicClientSchema(restConfig, group, version, resource), nil
	}

	return cluster
}
