package apis

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	resty "github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"

	"github.com/acornsoftlab/dashboard/pkg/app"
	"github.com/acornsoftlab/dashboard/pkg/config"
	"github.com/acornsoftlab/dashboard/pkg/lang"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	cmd "k8s.io/client-go/tools/clientcmd"
	cmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func ListContexts(c *gin.Context) {
	g := app.Gin{C: c}

	// query "ctx" 가 공백이면 default context  사용
	ctx := lang.NVL(c.Query("ctx"), config.Value.DefaultContext)

	resources := make(map[string]interface{})
	namespaces := []string{}

	if len(config.Value.Contexts) > 0 {

		// client
		kubeconfig, err := config.KubeConfigs(ctx)
		if err != nil {
			g.SendMessage(http.StatusBadRequest, err.Error(), err)
			return
		}

		// namespaces
		k8sClient, err := kubernetes.NewForConfig(kubeconfig)
		if err != nil {
			g.SendMessage(http.StatusInternalServerError, err.Error(), err)
			return
		}

		nsList, err := k8sClient.CoreV1().Namespaces().List(context.TODO(), v1.ListOptions{})
		if err != nil {
			g.SendMessage(http.StatusInternalServerError, err.Error(), err)
			return
		}
		for _, ns := range nsList.Items {
			namespaces = append(namespaces, ns.GetObjectMeta().GetName())
		}

		// resources
		discoveryClient, err := discovery.NewDiscoveryClientForConfig(kubeconfig)
		if err != nil {
			g.SendMessage(http.StatusInternalServerError, err.Error(), err)
			return
		}

		resourcesList, err := discoveryClient.ServerPreferredResources()
		if err != nil {
			g.SendMessage(http.StatusInternalServerError, err.Error(), err)
			return
		}

		// make a "groups > group > resources > resource" data structure
		var nm string
		for _, grpList := range resourcesList {
			if strings.Contains(grpList.GroupVersion, "/") {
				nm = strings.Split(grpList.GroupVersion, "/")[0]
			} else {
				nm = ""
			}

			if resources[nm] == nil {
				resources[nm] = make(map[string]interface{})
			}
			var group map[string]interface{}
			if resources[nm] == nil {
				group = make(map[string]interface{})
			} else {
				group = resources[nm].(map[string]interface{})
			}

			for _, r := range grpList.APIResources {
				group[r.Name] = map[string]interface{}{
					"name":         r.Name,
					"groupVersion": grpList.GroupVersion,
					"kind":         r.Kind,
					"namespaced":   r.Namespaced,
				}
			}

		}

	}

	g.Send(http.StatusOK, map[string]interface{}{
		"contexts": config.Value.Contexts,
		"currentContext": map[string]interface{}{
			"name":       ctx,
			"resources":  resources,
			"namespaces": namespaces,
		},
	})

}

func CreateContexts(c *gin.Context) {
	g := app.Gin{C: c}

	conf := cmdapi.Config{}

	if g.C.BindJSON(&conf) != nil {
		g.SendMessage(http.StatusBadRequest, "Unable to bind request body", nil)
	} else {

		err := cmd.WriteToFile(conf, config.Value.ConfigLoadingRules.GetDefaultFilename())
		if err != nil {
			g.SendMessage(http.StatusBadRequest, "Unable to modify kubeconfig", err)
		} else {
			config.Setup()
			reloadConfigMetricsScraper()
			ListContexts(g.C)
		}
	}
}

func GetContext(c *gin.Context) {
	g := app.Gin{C: c}

	conf := config.Value.KubeConfig.DeepCopy()
	cmdapi.ShortenConfig(conf)

	context := conf.Contexts[c.Param("CLUSTER")]
	if context == nil {
		g.SendMessage(http.StatusNotFound, fmt.Sprintf("not found a context '%s'", c.Param("CLUSTER")), nil)
	} else {
		g.Send(http.StatusOK, map[string]interface{}{
			"cluster": conf.Clusters[context.Cluster],
			"context": context,
			"user":    conf.AuthInfos[context.AuthInfo],
		})
	}
}

func CreateContext(c *gin.Context) {
	g := app.Gin{C: c}

	conf := config.Value.KubeConfig.DeepCopy()
	name := c.Param("CLUSTER")

	if conf != nil && conf.Contexts[name] != nil {
		g.SendMessage(http.StatusBadRequest, fmt.Sprintf("Already exist a context '%s'", name), nil)
	} else {
		json := make(map[string]interface{})
		if g.C.BindJSON(&json) != nil {
			g.SendMessage(http.StatusInternalServerError, "Unable to bind request body", nil)
			return
		}

		cluster := &cmdapi.Cluster{} // cluster
		context := &cmdapi.Context{} // context
		user := &cmdapi.AuthInfo{}   // user

		jsonC := json["cluster"].(map[string]interface{})
		jsonU := json["user"].(map[string]interface{})

		context.Cluster = fmt.Sprintf("%s-cluster", name)
		context.AuthInfo = fmt.Sprintf("%s-user", name)

		// cluster, user 중복 회피
		if config.Value.KubeConfig != nil {
			if _, exist := config.Value.KubeConfig.Clusters[context.Cluster]; exist {
				context.Cluster = fmt.Sprintf("%s-%s-cluster", name, lang.RandomString(5))
			}
			if _, exist := config.Value.KubeConfig.AuthInfos[context.AuthInfo]; exist {
				context.AuthInfo = fmt.Sprintf("%s-%s-user", name, lang.RandomString(5))
			}
		}

		log.Infof("%s, %s, %s", name, context.Cluster, context.AuthInfo)

		// cluster
		cluster.Server = jsonC["server"].(string)
		val, exists := jsonC["certificate-authority-data"].(string)
		if exists {
			ca, err := base64.StdEncoding.DecodeString(val)
			if err != nil {
				g.SendMessage(http.StatusBadRequest, "Unable to decode cerificate-authority-data", err)
				return
			}
			cluster.CertificateAuthorityData = ca
		} else {
			cluster.CertificateAuthority = jsonC["certificate-authority"].(string)
		}

		// user
		val, exists = jsonU["client-certificate-data"].(string)

		if exists {
			ca, err := base64.StdEncoding.DecodeString(val)
			if err != nil {
				g.SendMessage(http.StatusBadRequest, "Unable to decode client-certificate-data", err)
				return
			}
			user.ClientCertificateData = ca
		} else {
			user.ClientCertificate = jsonU["client-certificate"].(string)
		}

		val, exists = jsonU["client-key-data"].(string)
		if exists {
			ca, err := base64.StdEncoding.DecodeString(val)
			if err != nil {
				g.SendMessage(http.StatusBadRequest, "Unable to decode client-key-data", err)
				return
			}
			user.ClientKeyData = ca
		} else {
			user.ClientKey = jsonU["client-key"].(string)
		}

		//context
		if conf == nil {
			conf = &cmdapi.Config{
				Clusters:  make(map[string]*cmdapi.Cluster),
				AuthInfos: make(map[string]*cmdapi.AuthInfo),
				Contexts:  make(map[string]*cmdapi.Context),
			}
		}
		conf.Clusters[context.Cluster] = cluster
		conf.AuthInfos[context.AuthInfo] = user
		conf.Contexts[name] = context

		err := cmd.ModifyConfig(config.Value.ConfigLoadingRules, *conf, false)
		if err != nil {
			g.SendMessage(http.StatusBadRequest, "unable to modify kubeconfig", err)
		} else {
			reloadConfigMetricsScraper()
			config.Setup()
			ListContexts(g.C)
		}
	}
}

func DeleteContext(c *gin.Context) {
	g := app.Gin{C: c}

	conf := config.Value.KubeConfig.DeepCopy()
	name := c.Param("CLUSTER")

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

		err := cmd.ModifyConfig(config.Value.ConfigLoadingRules, *conf, false)
		if err != nil {
			g.SendMessage(http.StatusBadRequest, "Unable to modify kubeconfig", err)
		} else {
			reloadConfigMetricsScraper()
			config.Setup()
			ListContexts(g.C)
		}
	} else {
		g.SendMessage(http.StatusNotFound, fmt.Sprintf("not found context %s", name), nil)
	}
}

// metric-scraper config reload
func reloadConfigMetricsScraper() {
	client := resty.New()
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		Patch(fmt.Sprintf("%s/api/v1/config", *config.Value.MetricsScraperUrl))

	if err != nil {
		log.Errorf("Unable to metrics scraper config reload (cause=%v)", err)
	}
}
