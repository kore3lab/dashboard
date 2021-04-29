package _raw

/**
  참조
    https://kubernetes.io/docs/reference/using-api/api-concepts/
    https://github.com/gin-gonic/gin
*/

import (
	"k8s.io/apimachinery/pkg/conversion"
	"net/http"
	"net/url"
	"strings"

	"github.com/acornsoftlab/dashboard/pkg/app"
	"github.com/acornsoftlab/dashboard/pkg/client"
	"github.com/acornsoftlab/dashboard/pkg/config"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
)

// Get api group list
func GetAPIGroupList(c *gin.Context) {
	g := app.Gin{C: c}

	// kubeconfig
	conf, err := config.KubeConfigs(g.C.Param("CLUSTER"))
	if err != nil {
		g.SendMessage(http.StatusBadRequest, err.Error(), err)
		return
	}

	// instancing dynamic client
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(conf)
	if err != nil {
		g.SendError(err)
		return
	}

	groups, err := discoveryClient.ServerGroups()
	if err != nil {
		g.SendError(err)
		return
	}

	g.Send(http.StatusOK, groups)

}

// Create or Update
func ApplyRaw(c *gin.Context) {
	g := app.Gin{C: c}

	// kubeconfig
	conf, err := config.KubeConfigs(g.C.Param("CLUSTER"))
	if err != nil {
		g.SendMessage(http.StatusBadRequest, err.Error(), err)
		return
	}

	// api client
	api := client.NewDynamicClient(conf)

	// invoke POST
	r, err := api.POST(g.C.Request.Body, g.C.Request.Method == "PUT")
	if err != nil {
		g.SendMessage(http.StatusBadRequest, err.Error(), err)
		return
	}

	g.Send(http.StatusCreated, r)
}

// Delete
func DeleteRaw(c *gin.Context) {
	g := app.Gin{C: c}

	// url parameter validation
	v := []string{"VERSION", "RESOURCE", "NAME"}
	if err := g.ValidateUrl(v); err != nil {
		g.SendMessage(http.StatusBadRequest, err.Error(), err)
		return
	}

	// kubeconfig
	conf, err := config.KubeConfigs(g.C.Param("CLUSTER"))
	if err != nil {
		g.SendMessage(http.StatusBadRequest, err.Error(), err)
		return
	}

	// instancing dynamic client
	api := client.NewDynamicClientSchema(conf, c.Param("GROUP"), c.Param("VERSION"), c.Param("RESOURCE"))
	api.SetNamespace(c.Param("NAMESPACE"))

	// invoke delete
	if err := api.DELETE(c.Param("NAME"), v1.DeleteOptions{}); err != nil {
		g.SendError(err)
		return
	}

}

// Get or List
func GetRaw(c *gin.Context) {
	g := app.Gin{C: c}

	var err error

	// kubeconfig
	conf, err := config.KubeConfigs(g.C.Param("CLUSTER"))
	if err != nil {
		g.SendMessage(http.StatusBadRequest, err.Error(), err)
		return
	}

	var s conversion.Scope
	ListOptions := v1.ListOptions{}
	u, _ := url.Parse(c.Request.RequestURI)
	query, _ := url.ParseQuery(u.RawQuery)

	err = v1.Convert_url_Values_To_v1_ListOptions(&query,&ListOptions,s)
	if err != nil {
		g.SendMessage(http.StatusBadRequest, err.Error(), err)
		return
	}
	// instancing dynamic client
	api := client.NewDynamicClientSchema(conf, c.Param("GROUP"), c.Param("VERSION"), c.Param("RESOURCE"))
	api.SetNamespace(c.Param("NAMESPACE"))

	var r interface{}

	if c.Param("NAME") == "" {
		r, err = api.List(ListOptions)
		if err != nil {
			g.SendError(err)
			return
		}
	} else {
		r, err = api.GET(c.Param("NAME"), v1.GetOptions{})
		if err != nil {
			if strings.HasSuffix(err.Error(), "not found") {
				g.SendMessage(http.StatusNotFound, err.Error(), err)
			} else {
				g.SendError(err)
			}
			return
		}
	}

	g.Send(http.StatusOK, r)

}

// Patch
func PatchRaw(c *gin.Context) {
	g := app.Gin{C: c}

	var err error

	// url parameter validation
	v := []string{"VERSION", "RESOURCE", "NAME"}
	if err := g.ValidateUrl(v); err != nil {
		g.SendMessage(http.StatusBadRequest, err.Error(), err)
		return
	}

	// kubeconfig
	conf, err := config.KubeConfigs(g.C.Param("CLUSTER"))
	if err != nil {
		g.SendMessage(http.StatusBadRequest, err.Error(), err)
		return
	}

	// instancing dynamic client
	api := client.NewDynamicClientSchema(conf, c.Param("GROUP"), c.Param("VERSION"), c.Param("RESOURCE"))
	api.SetNamespace(c.Param("NAMESPACE"))

	var r interface{}

	r, err = api.PATCH(c.Param("NAME"), types.PatchType(c.ContentType()), c.Request.Body, v1.PatchOptions{})
	if err != nil {
		if strings.HasSuffix(err.Error(), "not found") {
			g.SendMessage(http.StatusNotFound, err.Error(), err)
		} else {
			g.SendError(err)
		}
		return
	}

	g.Send(http.StatusOK, r)

}
