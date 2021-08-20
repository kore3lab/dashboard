package apis

/**
  참조
    https://kubernetes.io/docs/reference/using-api/api-concepts/
    https://github.com/gin-gonic/gin
*/

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kore3lab/dashboard/pkg/app"
	"github.com/kore3lab/dashboard/pkg/config"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// Get api group list
func GetAPIGroupList(c *gin.Context) {
	g := app.Gin{C: c}

	// instancing dynamic client
	client, err := config.Cluster.Client(g.C.Param("CLUSTER"))
	if err != nil {
		g.SendMessage(http.StatusBadRequest, err.Error(), err)
		return
	}

	discoveryClient, err := client.NewDiscoveryClient()
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

	// api client
	client, err := config.Cluster.Client(g.C.Param("CLUSTER"))
	if err != nil {
		g.SendMessage(http.StatusBadRequest, err.Error(), err)
		return
	}

	api, err := client.NewDynamicClient()
	if err != nil {
		g.SendError(err)
		return
	}

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

	// instancing dynamic client
	client, err := config.Cluster.Client(g.C.Param("CLUSTER"))
	if err != nil {
		g.SendMessage(http.StatusBadRequest, err.Error(), err)
		return
	}

	api, err := client.NewDynamicClientSchema(c.Param("GROUP"), c.Param("VERSION"), c.Param("RESOURCE"))
	if err != nil {
		g.SendError(err)
		return
	}

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

	ListOptions := v1.ListOptions{}
	u, _ := url.Parse(c.Request.RequestURI)
	query, _ := url.ParseQuery(u.RawQuery)

	err = v1.Convert_url_Values_To_v1_ListOptions(&query, &ListOptions, nil)
	if err != nil {
		g.SendMessage(http.StatusBadRequest, err.Error(), err)
		return
	}
	// instancing dynamic client
	client, err := config.Cluster.Client(g.C.Param("CLUSTER"))
	if err != nil {
		g.SendMessage(http.StatusBadRequest, err.Error(), err)
		return
	}

	api, err := client.NewDynamicClientSchema(c.Param("GROUP"), c.Param("VERSION"), c.Param("RESOURCE"))
	if err != nil {
		g.SendError(err)
		return
	}

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

	// instancing dynamic client
	client, err := config.Cluster.Client(g.C.Param("CLUSTER"))
	if err != nil {
		g.SendMessage(http.StatusBadRequest, err.Error(), err)
		return
	}

	api, err := client.NewDynamicClientSchema(c.Param("GROUP"), c.Param("VERSION"), c.Param("RESOURCE"))
	if err != nil {
		g.SendError(err)
		return
	}
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
