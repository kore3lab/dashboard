package apis

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/kore3lab/dashboard/pkg/app"
	"github.com/kore3lab/dashboard/pkg/config"
	"github.com/kore3lab/dashboard/pkg/lang"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd/api"
)

func ListContexts(c *gin.Context) {
	g := app.Gin{C: c}

	g.Send(http.StatusOK, map[string]interface{}{
		"contexts":       config.Cluster.ClusterNames,
		"currentContext": lang.NVL(c.Query("ctx"), config.Cluster.DefaultContext),
	})

}

func GetContext(c *gin.Context) {
	g := app.Gin{C: c}

	// query "ctx" 가 공백이면 default context  사용
	ctx := lang.NVL(c.Param("CLUSTER"), config.Cluster.DefaultContext)

	resources := make(map[string]interface{})
	namespaces := []string{}
	kubernetesVersion := ""
	platform := ""
	if len(config.Cluster.ClusterNames) > 0 {

		// client
		client, err := config.Cluster.Client(ctx)
		if err != nil {
			g.SendMessage(http.StatusBadRequest, err.Error(), err)
			return
		}

		// namespaces
		k8sClient, err := client.NewKubernetesClient()
		if err != nil {
			g.SendMessage(http.StatusInternalServerError, err.Error(), err)
			return
		}

		timeout := int64(6)
		options := v1.ListOptions{TimeoutSeconds: &timeout} //timeout 6s
		nsList, err := k8sClient.CoreV1().Namespaces().List(context.TODO(), options)
		if err != nil {
			// namespace 를 가져오지 못하면 클라이언트에 contexts 리스트만 리턴해준다
			g.Send(http.StatusPartialContent, map[string]interface{}{
				"contexts": config.Cluster.ClusterNames,
				"currentContext": map[string]interface{}{
					"name":       ctx,
					"resources":  []string{},
					"namespaces": []string{},
				},
				"error": err.Error(),
			})
			return
		}
		for _, ns := range nsList.Items {
			namespaces = append(namespaces, ns.GetObjectMeta().GetName())
		}

		// resources
		discoveryClient, err := client.NewDiscoveryClient()
		if err != nil {
			g.SendMessage(http.StatusInternalServerError, err.Error(), err)
			return
		}

		ver, err := discoveryClient.ServerVersion()
		if err == nil {
			kubernetesVersion = ver.GitVersion
			platform = ver.Platform
		}

		resourcesList, err := discoveryClient.ServerPreferredResources()
		if err != nil {
			if _, resourcesList, err = discoveryClient.ServerGroupsAndResources(); err != nil {
				g.SendMessage(http.StatusInternalServerError, err.Error(), err)
				return
			}
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
		"contexts": config.Cluster.ClusterNames,
		"currentContext": map[string]interface{}{
			"name":              ctx,
			"resources":         resources,
			"namespaces":        namespaces,
			"kubernetesVersion": kubernetesVersion,
			"platform":          platform,
		},
	})

}

func GetContextNamespaces(c *gin.Context) {
	g := app.Gin{C: c}

	// query "ctx" 가 공백이면 default context  사용
	ctx := lang.NVL(c.Param("CLUSTER"), config.Cluster.DefaultContext)

	namespaces := []string{}

	// client
	client, err := config.Cluster.Client(ctx)
	if err != nil {
		g.SendMessage(http.StatusBadRequest, err.Error(), err)
		return
	}

	// namespaces
	k8sClient, err := client.NewKubernetesClient()
	if err != nil {
		g.SendMessage(http.StatusInternalServerError, err.Error(), err)
		return
	}

	timeout := int64(6)
	options := v1.ListOptions{TimeoutSeconds: &timeout} //timeout 6s
	nsList, err := k8sClient.CoreV1().Namespaces().List(context.TODO(), options)
	if err != nil {
		// namespace 를 가져오지 못하면 클라이언트에 contexts 리스트만 리턴해준다
		g.Send(http.StatusPartialContent, map[string]interface{}{
			"currentContext": ctx,
			"namespaces":     []string{},
		})
		return
	}
	for _, ns := range nsList.Items {
		namespaces = append(namespaces, ns.GetObjectMeta().GetName())
	}

	g.Send(http.StatusOK, map[string]interface{}{
		"currentContext": ctx,
		"namespaces":     namespaces,
	})

}
func DeleteContext(c *gin.Context) {
	g := app.Gin{C: c}

	if lang.ArrayContains(config.Cluster.ClusterNames, c.Param("CLUSTER")) {
		err := config.Cluster.Remove(c.Param("CLUSTER"))
		if err != nil {
			g.SendMessage(http.StatusBadRequest, "Unable to modify kubeconfig", err)
		} else {
			config.Setup()
			go reloadConfigMetricsScraper()
			go reloadConfigTerminal()
			ListContexts(c)
		}
	} else {
		g.SendMessage(http.StatusNotFound, fmt.Sprintf("Can't found a context '%s'", c.Param("CLUSTER")), nil)
	}

}

func GetContextConfig(c *gin.Context) {
	g := app.Gin{C: c}

	conf := config.Cluster.KubeConfig.DeepCopy()

	if g.C.Query("redacted") != "false" && g.C.Query("redacted") != "N" && g.C.Query("redacted") != "0" {
		api.ShortenConfig(conf)
	}

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

func CreateContexts(c *gin.Context) {
	g := app.Gin{C: c}

	conf := api.Config{}

	if g.C.BindJSON(&conf) != nil {
		g.SendMessage(http.StatusBadRequest, "Unable to bind request body", nil)
	} else {
		err := config.Cluster.Create(conf)
		if err != nil {
			g.SendMessage(http.StatusBadRequest, "Unable to modify kubeconfig", err)
		} else {
			config.Setup()
			go reloadConfigMetricsScraper()
			go reloadConfigTerminal()
			ListContexts(c)
		}
	}
}

func AddContext(c *gin.Context) {
	g := app.Gin{C: c}

	name := c.Param("CLUSTER")

	conf := config.Cluster.KubeConfig.DeepCopy()
	if conf != nil && conf.Contexts[name] != nil {
		g.SendMessage(http.StatusBadRequest, fmt.Sprintf("Already exist a context '%s'", name), nil)
	} else {
		json := make(map[string]interface{})
		if g.C.BindJSON(&json) != nil {
			g.SendMessage(http.StatusInternalServerError, "Unable to bind request body", nil)
			return
		}

		err := config.Cluster.Add(name, json)
		if err != nil {
			g.SendMessage(http.StatusBadRequest, "Unable to modify kubeconfig", err)
		} else {
			config.Setup()
			go reloadConfigMetricsScraper()
			go reloadConfigTerminal()
			ListContexts(c)
		}
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

// terminal config reload
func reloadConfigTerminal() {
	client := resty.New()
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		Patch(fmt.Sprintf("%s/api/v1/config", config.Value.TerminalUrl))

	if err != nil {
		log.Errorf("Unable to Terminal config reload (cause=%v)", err)
	}
}
