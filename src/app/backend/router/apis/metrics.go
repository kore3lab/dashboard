package apis

import (
	//"errors"
	//"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kore3lab/dashboard/model"
	"github.com/kore3lab/dashboard/pkg/app"
	"github.com/kore3lab/dashboard/pkg/config"
	"github.com/kore3lab/dashboard/pkg/lang"
)

// Get node metrics
func GetClusterMetrics(c *gin.Context) {
	g := app.Gin{C: c}

	cluster := lang.NVL(g.C.Param("CLUSTER"), config.Cluster.DefaultContext)

	metrics, err := model.GetClusterCumulativeMetrics(cluster)
	if err != nil {
		g.SendError(err)
	} else {
		g.Send(http.StatusOK, metrics)
	}

}

// Get node metrics
func GetNodeMetrics(c *gin.Context) {
	g := app.Gin{C: c}

	cluster := lang.NVL(g.C.Param("CLUSTER"), config.Cluster.DefaultContext)

	metrics, err := model.GetNodeCumulativeMetrics(cluster, c.Param("NAME"))
	if err != nil {
		g.SendError(err)
	} else {
		g.Send(http.StatusOK, metrics)
	}

}

// Get workload metrics (pod, deployments, statefulsets, daemonsets, replicasets)
func GetWorkloadMetrics(c *gin.Context) {
	g := app.Gin{C: c}

	cluster := lang.NVL(g.C.Param("CLUSTER"), config.Cluster.DefaultContext)

	metrics, err := model.GetWorkloadCumulativeMetrics(cluster, c.Param("NAMESPACE"), c.Param("RESOURCE"), c.Param("NAME"))
	if err != nil {
		g.SendError(err)
	} else {
		g.Send(http.StatusOK, metrics)
	}

}

// Get node list
func GetNodeListWithUsage(c *gin.Context) {
	g := app.Gin{C: c}
	cluster := lang.NVL(g.C.Param("CLUSTER"), config.Cluster.DefaultContext)

	pods, err := model.GetNodeListWithUsage(cluster)
	if err != nil {
		g.SendError(err)
	} else {
		g.Send(http.StatusOK, pods)
	}

}

// Get node pod-list
func GetNodePodListWithMetrics(c *gin.Context) {
	g := app.Gin{C: c}
	cluster := lang.NVL(g.C.Param("CLUSTER"), config.Cluster.DefaultContext)

	pods, err := model.GetNodePodListWithMetrics(cluster, c.Param("NAME"))
	if err != nil {
		g.SendError(err)
	} else {
		g.Send(http.StatusOK, pods)
	}

}

// Get workload pod-list (deployments, statefulsets, daemonsets, replicasets)
func GetWorkloadPodListWithMetrics(c *gin.Context) {
	g := app.Gin{C: c}
	cluster := lang.NVL(g.C.Param("CLUSTER"), config.Cluster.DefaultContext)

	pods, err := model.GetWorkloadPodListWithMetrics(cluster, c.Param("NAMESPACE"), c.Param("RESOURCE"), c.Param("NAME"))
	if err != nil {
		g.SendError(err)
	} else {
		g.Send(http.StatusOK, pods)
	}

}
