package apis

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/acornsoftlab/dashboard/pkg/app"
	"github.com/acornsoftlab/dashboard/pkg/config"
	"github.com/gin-gonic/gin"
	cmd "k8s.io/client-go/tools/clientcmd"
	cmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func ListContexts(c *gin.Context) {
	g := app.Gin{C: c}

	g.Send(http.StatusOK, map[string]interface{}{
		"currentContext": config.Value.CurrentContext,
		"contexts":       config.Value.Contexts,
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
		g.SendMessage(http.StatusBadRequest, fmt.Sprintf("exist a context '%s'", name))
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
			config.Setup()
			ListContexts(g.C)
		}
	}
}

func DeleteContext(c *gin.Context) {
	g := app.Gin{C: c}

	conf := config.Value.KubeConfig.DeepCopy()
	name := c.Param("CLUSTER")

	fmt.Println(conf.Contexts[name])

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
			config.Setup()
			ListContexts(g.C)
		}
	} else {
		g.SendMessage(http.StatusNotFound, fmt.Sprintf("not found context %s", name))
	}
}
