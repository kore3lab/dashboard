package model

// base code : https://github.com/kubernetes/dashboard/tree/master/src/app/backend/resource/daemonset

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/kore3lab/dashboard/pkg/config"
	"github.com/kore3lab/dashboard/pkg/lang"
	log "github.com/sirupsen/logrus"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	metricsapi "k8s.io/metrics/pkg/apis/metrics"
	metricsV1beta1api "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

type nodeMetric struct {
	CPU     int64 `json:"cpu"`
	Memory  int64 `json:"memory"`
	Storage int64 `json:"storage"`
	Pods    int64 `json:"pods"`
}
type nodeRatio struct {
	CPU     float32 `json:"cpu"`
	Memory  float32 `json:"memory"`
	Storage float32 `json:"storage"`
	Pods    float32 `json:"pods"`
}

type NodeWithMetrics struct {
	Roles             string           `json:"role"`
	Address           string           `json:"address"`
	Version           string           `json:"version"`
	Metrics           NodeMetricsUsage `json:"metrics"`
	CreationTimestamp metaV1.Time      `json:"creationTimestamp"`
	Status            string           `json:"status"`
}

type NodeMetricsUsage struct {
	Allocatable nodeMetric `json:"allocatable"`
	Usage       nodeMetric `json:"usage"`
	Percent     nodeRatio  `json:"percent"`
}

// returns a subset of pods controlled by given deployment.
func GetNodePods(apiClient *kubernetes.Clientset, name string) ([]coreV1.Pod, error) {

	podList, err := apiClient.CoreV1().Pods("").List(context.TODO(), metaV1.ListOptions{FieldSelector: fmt.Sprintf("spec.nodeName=%s", name)})
	if err != nil {
		return nil, err
	}
	return podList.Items, nil

}

// get node-list with metrics-usage
func GetNodeListWithUsage(cluster string) (interface{}, error) {

	clientSet, err := config.Cluster.Client(cluster)
	if err != nil {
		return nil, err
	}

	apiClient, err := clientSet.NewKubernetesClient()
	if err != nil {
		return nil, err
	}

	metricsClient, err := clientSet.NewMetricsClient()
	if err != nil {
		return nil, err
	}

	//timeout 5s
	timeout := int64(5)

	// node-list
	nodeList, err := apiClient.CoreV1().Nodes().List(context.TODO(), metaV1.ListOptions{TimeoutSeconds: &timeout})
	if err != nil {
		return nil, err
	}

	// self.Workloads.Pods (노드별 파드 수 & running 파드 수)
	podList, err := apiClient.CoreV1().Pods("").List(context.TODO(), metaV1.ListOptions{TimeoutSeconds: &timeout})
	if err != nil {
		return nil, err
	}
	usagePods := map[string]int64{}
	for _, m := range podList.Items {
		if m.Spec.NodeName != "" {
			usagePods[m.Spec.NodeName] = usagePods[m.Spec.NodeName] + 1
		}
	}

	nodes := map[string]*NodeWithMetrics{}
	summary := NodeMetricsUsage{}
	d := time.Duration(timeout) * time.Second
	nodeSummary := ProxyNodeSummary{}

	for _, m := range nodeList.Items {

		// node summary for storage used percentage (/api/v1/nodes/<node name>/proxy/stats/summary)
		request := apiClient.CoreV1().RESTClient().Get().Resource("nodes").Name(m.Name).SubResource("proxy").Suffix("stats/summary").Timeout(d)
		responseRawArrayOfBytes, err := request.DoRaw(context.Background())
		if err != nil {
			log.Warnf("Unable to get %s/proxy/stats/summary (cause=%v)", m, err)
		} else {
			err = json.Unmarshal(responseRawArrayOfBytes, &nodeSummary)
			if err != nil {
				log.Warnf("Unable to unmarshal data %s/proxy/stats/summary (cause=%v)", m, err)
			}
		}

		// allocatable <a node>
		nodes[m.Name] = &NodeWithMetrics{
			Address:           m.Status.Addresses[0].Address,
			Version:           m.Status.NodeInfo.KubeletVersion,
			CreationTimestamp: m.GetCreationTimestamp(),
			Status:            findNodeStatus(m),
			Roles:             findNodeRoles(m),
			Metrics: NodeMetricsUsage{
				Allocatable: nodeMetric{
					CPU:     m.Status.Allocatable.Cpu().MilliValue(),
					Memory:  m.Status.Allocatable.Memory().Value(),
					Storage: nodeSummary.Node.Fs.Capacitybytes,
					Pods:    m.Status.Allocatable.Pods().Value(),
				},
				Usage: nodeMetric{
					Storage: nodeSummary.Node.Fs.Usedbytes,
					Pods:    usagePods[m.Name],
				},
				Percent: nodeRatio{
					Storage: float32(lang.DivideRound(nodeSummary.Node.Fs.Usedbytes, nodeSummary.Node.Fs.Capacitybytes, 4) * 100),
					Pods:    float32(lang.DivideRound(usagePods[m.Name], m.Status.Allocatable.Pods().Value(), 4) * 100),
				},
			},
		}

		// (+) allocate/usage <total nodes>
		summary.Allocatable.CPU += m.Status.Allocatable.Cpu().MilliValue()
		summary.Allocatable.Memory += m.Status.Allocatable.Memory().Value()
		summary.Allocatable.Pods += m.Status.Allocatable.Pods().Value()
		summary.Allocatable.Storage += nodeSummary.Node.Fs.Capacitybytes
		summary.Usage.Storage += nodeSummary.Node.Fs.Usedbytes
		summary.Usage.Pods += usagePods[m.Name]

	}

	// get metrics - (cpu, memory)
	versionedNodeMetrics, err := metricsClient.MetricsV1beta1().NodeMetricses().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		log.Warnf("Unable to get node metries (cause=%v)", err)
	}

	nodeMetrics := &metricsapi.NodeMetricsList{}
	err = metricsV1beta1api.Convert_v1beta1_NodeMetricsList_To_metrics_NodeMetricsList(versionedNodeMetrics, nodeMetrics, nil)
	if err != nil {
		log.Warnf("Unable to convert node metries (cause=%v)", err)
	}

	for _, m := range nodeMetrics.Items {
		// usage <a node>
		nodes[m.Name].Metrics.Usage.CPU = m.Usage.Cpu().MilliValue()
		nodes[m.Name].Metrics.Usage.Memory = m.Usage.Memory().Value()
		nodes[m.Name].Metrics.Percent.CPU = float32(lang.DivideRound(m.Usage.Cpu().MilliValue(), nodes[m.Name].Metrics.Allocatable.CPU, 4) * 100)
		nodes[m.Name].Metrics.Percent.Memory = float32(lang.DivideRound(m.Usage.Memory().Value(), nodes[m.Name].Metrics.Allocatable.Memory, 4) * 100)

		// (+) usage <total nodes>
		summary.Usage.CPU += m.Usage.Cpu().MilliValue()
		summary.Usage.Memory += m.Usage.Memory().Value()

	}

	summary.Percent = nodeRatio{
		CPU:     float32(lang.DivideRound(summary.Usage.CPU, summary.Allocatable.CPU, 4) * 100),
		Memory:  float32(lang.DivideRound(summary.Usage.Memory, summary.Allocatable.Memory, 4) * 100),
		Storage: float32(lang.DivideRound(summary.Usage.Storage, summary.Allocatable.Storage, 4) * 100),
		Pods:    float32(lang.DivideRound(summary.Usage.Pods, summary.Allocatable.Pods, 4) * 100),
	}

	return struct {
		Nodes   map[string]*NodeWithMetrics `json:"nodes"`
		Summary NodeMetricsUsage            `json:"summary"`
	}{
		Nodes:   nodes,
		Summary: summary,
	}, nil

}

func findNodeStatus(node coreV1.Node) string {

	for _, c := range node.Status.Conditions {
		if c.Type == coreV1.NodeReady {
			if c.Status == coreV1.ConditionTrue {
				return "Ready"
			} else {
				return "NotReady"
			}
		}
	}
	return "Unknown"
}

// findNodeRoles returns the roles of a given node.
// The roles are determined by looking for:
// * a node-role.kubernetes.io/<role>="" label
// * a kubernetes.io/role="<role>" label
func findNodeRoles(node coreV1.Node) string {
	roleList := sets.NewString()
	for k, v := range node.Labels {
		switch {
		case strings.HasPrefix(k, LabelNodeRolePrefix):
			if role := strings.TrimPrefix(k, LabelNodeRolePrefix); len(role) > 0 {
				roleList.Insert(role)
			}

		case k == NodeLabelRole && v != "":
			roleList.Insert(v)
		}
	}

	if len(roleList) > 0 {
		return strings.Join(roleList.List(), ",")
	} else {
		return "<none>"
	}

}

const (
	LabelNodeRolePrefix = "node-role.kubernetes.io/"
	NodeLabelRole       = "kubernetes.io/role"
)
