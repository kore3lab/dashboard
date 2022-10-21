package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kore3lab/dashboard/model"
	"github.com/kore3lab/dashboard/pkg/app"
	"github.com/kore3lab/dashboard/pkg/config"
	"github.com/kore3lab/dashboard/pkg/lang"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Network(c *gin.Context) {
	g := app.Gin{C: c}

	cluster := lang.NVL(g.C.Param("CLUSTER"), config.Cluster.DefaultContext)
	namespace := c.Param("NAMESPACE")

	if topology, err := model.GetNetworkGraph(cluster, namespace); err != nil {
		g.SendError(err)
	} else {
		g.Send(http.StatusOK, topology)
	}

}

func Topology(c *gin.Context) {
	g := app.Gin{C: c}

	cluster := lang.NVL(g.C.Param("CLUSTER"), config.Cluster.DefaultContext)
	namespace := c.Param("NAMESPACE")

	if topology, err := model.GetTopologyGraph(cluster, namespace); err != nil {
		g.SendError(err)
	} else {
		g.Send(http.StatusOK, topology)
	}

}

func Workloads(c *gin.Context) {
	g := app.Gin{C: c}

	cluster := lang.NVL(g.C.Param("CLUSTER"), config.Cluster.DefaultContext)
	namespace := c.Param("NAMESPACE")

	if workloads, err := model.GetWorkloadGraph(cluster, namespace); err != nil {
		g.SendError(err)
	} else {
		g.Send(http.StatusOK, workloads)
	}

}

func Pod(c *gin.Context) {
	g := app.Gin{C: c}

	cluster := lang.NVL(g.C.Param("CLUSTER"), config.Cluster.DefaultContext)
	namespace := c.Param("NAMESPACE")
	name := c.Param("POD")

	if workloads, err := model.GetPodGraph(cluster, namespace, name); err != nil {
		g.SendError(err)
	} else {
		g.Send(http.StatusOK, workloads)
	}

}

func Dashboard(c *gin.Context) {
	g := app.Gin{C: c}

	cluster := lang.NVL(g.C.Param("CLUSTER"), config.Cluster.DefaultContext)

	clientSet, err := config.Cluster.Client(cluster)
	if err != nil {
		g.SendError(err)
		return
	}

	apiClient, err := clientSet.NewKubernetesClient()
	if err != nil {
		g.SendError(err)
		return
	}

	timeout := int64(5)
	options := metaV1.ListOptions{TimeoutSeconds: &timeout} //timeout 5s

	workloads := map[string]struct {
		Ready     int `json:"ready"`
		Available int `json:"available"`
	}{}

	// daemonset
	available, ready, _ := model.GetDaemonSetsReady(apiClient, options)
	workloads["daemonset"] = struct {
		Ready     int `json:"ready"`
		Available int `json:"available"`
	}{Available: available, Ready: ready}

	// deployment
	available, ready, _ = model.GetDeploymentsReady(apiClient, options)
	workloads["deployment"] = struct {
		Ready     int `json:"ready"`
		Available int `json:"available"`
	}{Available: available, Ready: ready}

	// replicaset
	available, ready, _ = model.GetReplicaSetsReady(apiClient, options)
	workloads["replicaset"] = struct {
		Ready     int `json:"ready"`
		Available int `json:"available"`
	}{Available: available, Ready: ready}

	// statefulset
	available, ready, _ = model.GetStatefulSetsReady(apiClient, options)
	workloads["statefulset"] = struct {
		Ready     int `json:"ready"`
		Available int `json:"available"`
	}{Available: available, Ready: ready}

	// pods
	available, ready, _ = model.GetPodsReady(apiClient, options)
	workloads["pod"] = struct {
		Ready     int `json:"ready"`
		Available int `json:"available"`
	}{Available: available, Ready: ready}

	g.Send(http.StatusOK, workloads)
}
