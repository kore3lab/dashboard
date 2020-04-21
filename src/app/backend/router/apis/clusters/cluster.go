package apis

import (
	"net/http"

	"github.com/acornsoftlab/kore3/model/v1alpha1"
	"github.com/acornsoftlab/kore3/pkg/app"
	"github.com/gin-gonic/gin"
)

func Topology(c *gin.Context) {
	g := app.Gin{C: c}

	// parameter validation
	if err := g.ValidateUrl([]string{"CLUSTER"}); err != nil {
		g.SendMessage(http.StatusBadRequest, err.Error())
		return
	}
	cluster := c.Param("CLUSTER")
	namespace := c.Param("NAMESPACE")

	topology := model.NewTopology(cluster)
	if err := topology.Get(namespace); err != nil {
		g.Send(500, err.Error())
	}

	g.Send(http.StatusOK, topology)
	// g.SendMessage(http.StatusOK, fmt.Sprintf("There are %d pods in the cluster %s", len(pods.Items), cluster))

}

func Dashboard(c *gin.Context) {
	g := app.Gin{C: c}

	// parameter validation
	if err := g.ValidateUrl([]string{"CLUSTER"}); err != nil {
		g.SendMessage(http.StatusBadRequest, err.Error())
		return
	}
	cluster := c.Param("CLUSTER")

	dashboard := model.NewDashboard(cluster)
	if err := dashboard.Get(); err != nil {
		g.Send(500, err.Error())
	}

	g.Send(http.StatusOK, dashboard)
	// g.SendMessage(http.StatusOK, fmt.Sprintf("There are %d pods in the cluster %s", len(pods.Items), cluster))

}
