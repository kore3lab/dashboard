package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kore3lab/dashboard/model"
	"github.com/kore3lab/dashboard/pkg/app"
	"github.com/kore3lab/dashboard/pkg/config"
	"github.com/kore3lab/dashboard/pkg/lang"
)

func Topology(c *gin.Context) {
	g := app.Gin{C: c}

	cluster := lang.NVL(g.C.Param("CLUSTER"), config.Cluster.DefaultContext)
	namespace := c.Param("NAMESPACE")

	topology := model.NewTopology(cluster)
	if err := topology.Get(namespace); err != nil {
		g.SendError(err)
	} else {
		g.Send(http.StatusOK, topology)
	}

}

func Dashboard(c *gin.Context) {
	g := app.Gin{C: c}

	cluster := lang.NVL(g.C.Param("CLUSTER"), config.Cluster.DefaultContext)

	dashboard := model.NewDashboard(cluster)
	if err := dashboard.Get(); err != nil {
		g.SendError(err)
	} else {
		g.Send(http.StatusOK, dashboard)
	}

}
