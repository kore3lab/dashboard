package apis

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kore3lab/dashboard/model/v1alpha1"
	"github.com/kore3lab/dashboard/pkg/app"
	"github.com/kore3lab/dashboard/pkg/config"
	"github.com/kore3lab/dashboard/pkg/lang"
)

// Get node metrics
func GetNodeMetrics(c *gin.Context) {
	g := app.Gin{C: c}

	cluster := lang.NVL(g.C.Param("CLUSTER"), config.Cluster.DefaultContext)

	scraperClient, err := model.NewScraperClient(cluster)
	if err != nil {
		g.SendMessage(http.StatusInternalServerError, "Unable to get scrapping client", err)
		return
	}

	err = scraperClient.GetNodeMetrics(c.Param("NODE"))
	if err != nil {
		g.SendMessage(http.StatusInternalServerError, "Unable to get scrapping metrics", err)
	} else {
		g.Send(http.StatusOK, scraperClient)
	}

}

// Get metrics (pod, deployments, statefulsets, daemonsets, replicasets)
func GetMetrics(c *gin.Context) {
	g := app.Gin{C: c}

	cluster := lang.NVL(g.C.Param("CLUSTER"), config.Cluster.DefaultContext)

	scraperClient, err := model.NewScraperClient(cluster)
	if err != nil {
		g.SendMessage(http.StatusInternalServerError, "Unable to get scrapping client", err)
		return
	}

	if c.Param("RESOURCE") == "pods" {
		err = scraperClient.GetNodeMetrics(c.Param("NAME"))
	} else if c.Param("RESOURCE") == "deployments" {
		err = scraperClient.GetDeploymentMetrics(c.Param("NAMESPACE"), c.Param("NAME"))
	} else if c.Param("RESOURCE") == "statefulsets" {
		err = scraperClient.GetSetatefulSetMetrics(c.Param("NAMESPACE"), c.Param("NAME"))
	} else if c.Param("RESOURCE") == "daemonsets" {
		err = scraperClient.GetDaemonSetMetrics(c.Param("NAMESPACE"), c.Param("NAME"))
	} else if c.Param("RESOURCE") == "replicasets" {
		err = scraperClient.GetReplicaSetMetrics(c.Param("NAMESPACE"), c.Param("NAME"))
	} else {
		err = errors.New(fmt.Sprintf("unsupported resource '%s'", c.Param("RESOURCE")))
	}

	if err != nil {
		g.SendMessage(http.StatusInternalServerError, "Unable to get scrapping metrics", err)
	} else {
		g.Send(http.StatusOK, scraperClient)
	}

}
