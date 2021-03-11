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

	//ctx parameter 가 CurrentContext와 다르고 Contexts 에 포함되어 있다면 CurrentContext 변경
	ctx := c.DefaultQuery("ctx", config.Value.CurrentContext)
	if ctx != config.Value.CurrentContext && lang.ArrayContains(config.Value.Contexts, ctx) {
		config.Value.CurrentContext = ctx
	}

	// client
	kubeconfig := config.Value.KubeConfigs[config.Value.CurrentContext]

	// namespaces
	k8sClient, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		g.SendMessage(http.StatusInternalServerError, err.Error())
		return
	}

	nsList, err := k8sClient.CoreV1().Namespaces().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		g.SendMessage(http.StatusInternalServerError, err.Error())
		return
	}
	namespaces := []string{}
	for _, ns := range nsList.Items {
		namespaces = append(namespaces, ns.GetObjectMeta().GetName())
	}

	// resources
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(kubeconfig)
	if err != nil {
		g.SendMessage(http.StatusInternalServerError, err.Error())
		return
	}

	resourcesList, err := discoveryClient.ServerPreferredResources()
	if err != nil {
		g.SendMessage(http.StatusInternalServerError, err.Error())
		return
	}

	// make a "groups > group > resources > resource" data structure
	var nm string
	resources := make(map[string]interface{})
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

	g.Send(http.StatusOK, map[string]interface{}{
		"contexts": config.Value.Contexts,
		"currentContext": map[string]interface{}{
			"name":       config.Value.CurrentContext,
			"resources":  resources,
			"namespaces": namespaces,
		},
	})

}

func CreateContexts(c *gin.Context) {
	g := app.Gin{C: c}

	conf := cmdapi.Config{}

	if g.C.BindJSON(&conf) != nil {
		g.SendMessage(http.StatusInternalServerError, "Unable to bind request body")
	} else {
		err := cmd.ModifyConfig(config.Value.ConfigLoadingRules, conf, false)
		if err != nil {
			g.SendMessage(http.StatusBadRequest, "Unable to modify kubeconfig")
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
		g.SendMessage(http.StatusNotFound, fmt.Sprintf("not found a context '%s'", c.Param("CLUSTER")))
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

	if conf.Contexts[name] != nil {
		g.SendMessage(http.StatusBadRequest, fmt.Sprintf("Already exist a context '%s'", name))
	} else {
		json := make(map[string]interface{})
		if g.C.BindJSON(&json) != nil {
			g.SendMessage(http.StatusInternalServerError, "Unable to bind request body")
			return
		}

		cluster := &cmdapi.Cluster{}
		context := &cmdapi.Context{}
		user := &cmdapi.AuthInfo{}

		jsonC := json["cluster"].(map[string]interface{})
		jsonU := json["user"].(map[string]interface{})

		context.Cluster = name
		context.AuthInfo = fmt.Sprintf("%s", name)

		cluster.Server = jsonC["server"].(string)
		val, exists := jsonC["certificate-authority-data"].(string)
		if exists {
			ca, err := base64.StdEncoding.DecodeString(val)
			if err != nil {
				g.SendMessage(http.StatusBadRequest, "Unable to decode cerificate-authority-data")
				return
			}
			cluster.CertificateAuthorityData = ca
		} else {
			cluster.CertificateAuthority = jsonC["certificate-authority"].(string)
		}

		val, exists = jsonU["client-certificate-data"].(string)

		if exists {
			ca, err := base64.StdEncoding.DecodeString(val)
			if err != nil {
				g.SendMessage(http.StatusBadRequest, "Unable to decode client-certificate-data")
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
				g.SendMessage(http.StatusBadRequest, "Unable to decode client-key-data")
				return
			}
			user.ClientKeyData = ca
		} else {
			user.ClientKey = jsonU["client-key"].(string)
		}

		conf.Clusters[context.Cluster] = cluster
		conf.AuthInfos[context.AuthInfo] = user
		conf.Contexts[name] = context

		err := cmd.ModifyConfig(config.Value.ConfigLoadingRules, *conf, false)
		if err != nil {
			g.SendMessage(http.StatusBadRequest, "unable to modify kubeconfig")
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
			g.SendMessage(http.StatusBadRequest, "Unable to modify kubeconfig")
		} else {
			reloadConfigMetricsScraper()
			config.Setup()
			ListContexts(g.C)
		}
	} else {
		g.SendMessage(http.StatusNotFound, fmt.Sprintf("not found context %s", name))
	}
}

// metric-scraper config reload
func reloadConfigMetricsScraper() {
	client := resty.New()
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		Patch(fmt.Sprintf("%s/api/v1/config", config.Value.MetricsScraperUrl))

	if err != nil {
		log.Errorf("Unable to metrics scraper config reload (cause=%v)", err)
	}
}
