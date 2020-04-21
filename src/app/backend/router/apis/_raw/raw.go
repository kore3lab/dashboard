package _raw

/**
  참조
    https://kubernetes.io/docs/reference/using-api/api-concepts/
    https://github.com/gin-gonic/gin
*/

import (
	"net/http"

	"github.com/acornsoftlab/kore3/pkg/app"
	"github.com/acornsoftlab/kore3/pkg/client"
	"github.com/gin-gonic/gin"

	"github.com/acornsoftlab/kore3/pkg/config"
)

// Create or Update
func ApplyRaw(c *gin.Context) {
	g := app.Gin{C: c}

	// url parameter validation
	if err := g.ValidateUrl([]string{"CLUSTER"}); err != nil {
		g.SendMessage(http.StatusBadRequest, err.Error())
		return
	}

	// api clinet
	context := g.C.Param("CLUSTER")
	api := client.NewDynamicClient(config.Value.KubeConfigs[context])

	// invoke POST
	reader := g.C.Request.Body
	defer reader.Close()

	r, err := api.POST(g.C.Request.Body, g.C.Request.Method == "PUT")
	if err != nil {
		g.SendMessage(http.StatusBadRequest, err.Error())
		return
	}

	g.Send(http.StatusCreated, r)
}

// Delete
func DeleteRaw(c *gin.Context) {
	g := app.Gin{C: c}

	// url parameter validation
	v := []string{"CLUSTER", "GROUP", "VERSION", "RESOURCETYPE", "NAME"}
	if err := g.ValidateUrl(v); err != nil {
		g.SendMessage(http.StatusBadRequest, err.Error())
		return
	}

	// api clinet
	context := g.C.Param("CLUSTER")
	api := client.NewDynamicClient(config.Value.KubeConfigs[context])

	// invoke DELETE
	namespace := c.Param("NAMESPACE")
	namespaceSet := (namespace != "")
	group := c.Param("GROUP")
	version := c.Param("VERSION")
	kind := c.Param("RESOURCETYPE")
	name := c.Param("NAME")
	if err := api.DELETE(namespaceSet, namespace, group, version, kind, name); err != nil {
		g.SendMessage(http.StatusInternalServerError, err.Error())
		return
	}

}

// Get or List
func GetRaw(c *gin.Context) {
	g := app.Gin{C: c}

	// url parameter validation
	v := []string{"CLUSTER", "VERSION", "RESOURCETYPE"}
	if err := g.ValidateUrl(v); err != nil {
		g.SendMessage(http.StatusBadRequest, err.Error())
		return
	}

	// api clinet
	context := g.C.Param("CLUSTER")
	api := client.NewDynamicClient(config.Value.KubeConfigs[context])

	// invoke DELETE
	namespace := c.Param("NAMESPACE")
	namespaceSet := (namespace != "")
	group := c.Param("GROUP")
	version := c.Param("VERSION")
	kind := c.Param("RESOURCETYPE")
	name := c.Param("NAME")

	var r interface{}
	var err error

	if name == "" {
		r, err = api.List(namespaceSet, namespace, group, version, kind)
		if err != nil {
			g.SendMessage(http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		r, err = api.GET(namespaceSet, namespace, group, version, kind, name)
		if err != nil {
			g.SendMessage(http.StatusInternalServerError, err.Error())
			return
		}
	}

	g.Send(http.StatusOK, r)

}
