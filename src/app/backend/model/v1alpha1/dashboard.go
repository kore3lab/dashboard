package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/acornsoftlab/kore3/pkg/config"
	"github.com/acornsoftlab/kore3/pkg/lang"
	"github.com/dustin/go-humanize"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	metricsapi "k8s.io/metrics/pkg/apis/metrics"
	metricsV1beta1api "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
)

type Dashboard struct {
	context   string
	Cluster   dashboardCluster          `json:"cluster"`
	Nodes     map[string]*dashboardNode `json:"nodes"`
	Workloads dashboardWorkload         `json:"workloads"`
}
type dashboardCluster struct {
	Nodes   string `json:"nodes"`
	Cpu     string `json:"cpu"`
	Memory  string `json:"memory"`
	Storage string `json:"storage"`
	Pods    string `json:"pods"`
}
type dashboardNode struct {
	Address string                 `json:"address"`
	Status  string                 `json:"status"`
	Roles   string                 `json:"roles"`
	Cpu     *dashboardNodeResource `json:"cpu"`
	Memory  *dashboardNodeResource `json:"memory"`
	Storage *dashboardNodeResource `json:"storage"`
	Pods    *dashboardNodeResource `json:"pods"`
}
type dashboardWorkload struct {
	DaemonSet   dashboardAvailable `json:"daemonset"`
	Deployment  dashboardAvailable `json:"deployment"`
	ReplicaSet  dashboardAvailable `json:"replicaset"`
	StatefulSet dashboardAvailable `json:"statefulset"`
	Pod         dashboardAvailable `json:"pod"`
}
type dashboardAvailable struct {
	Ready     int `json:"ready"`
	Available int `json:"available"`
}

type dashboardNodeResource struct {
	Allocatable string `json:"allocatable"`
	Usage       string `json:"usage"`
	Percent     string `json:"percent"`
}
type resource struct {
	cpu     int64
	memory  int64
	storage int64
	pods    int64
}

func NewDashboard(contextName string) Dashboard {
	return Dashboard{
		context: contextName,
		Nodes:   make(map[string]*dashboardNode),
	}
}

func (self *Dashboard) Get() error {

	conf := config.Value.KubeConfigs[self.context]

	//api client
	apiClient, err := kubernetes.NewForConfig(conf)
	if err != nil {
		return err
	}

	//metrics client
	metricsClient, err := metricsclientset.NewForConfig(conf)
	if err != nil {
		return err
	}

	// node list
	nodeList, err := apiClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	// pod list
	podList, err := apiClient.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	// daemonset list (ready count)
	dsList, err := apiClient.AppsV1().DaemonSets("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	ready := 0
	for _, m := range dsList.Items {
		if m.Status.NumberAvailable == m.Status.NumberAvailable {
			ready += 1
		}
	}
	self.Workloads.DaemonSet = dashboardAvailable{Available: len(dsList.Items), Ready: ready}

	// deployment list (ready count)
	deployList, err := apiClient.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	ready = 0
	for _, m := range deployList.Items {
		if m.Status.AvailableReplicas == m.Status.ReadyReplicas {
			ready += 1
		}
	}
	self.Workloads.Deployment = dashboardAvailable{Available: len(deployList.Items), Ready: ready}

	// replicaSet list (ready count)
	rsList, err := apiClient.AppsV1().ReplicaSets("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	ready = 0
	for _, m := range rsList.Items {
		if m.Status.AvailableReplicas == m.Status.ReadyReplicas {
			ready += 1
		}
	}
	self.Workloads.ReplicaSet = dashboardAvailable{Available: len(rsList.Items), Ready: ready}

	// setatefulset list (ready count)
	sfsList, err := apiClient.AppsV1().StatefulSets("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	ready = 0
	for _, m := range sfsList.Items {
		if m.Status.Replicas == m.Status.ReadyReplicas {
			ready += 1
		}
	}
	self.Workloads.StatefulSet = dashboardAvailable{Available: len(sfsList.Items), Ready: ready}

	// 노드 status, 리소스 allocatable
	allocate := map[string]resource{}
	allocateTotal := resource{}

	for _, m := range nodeList.Items {
		r := resource{
			cpu:     m.Status.Allocatable.Cpu().MilliValue(),
			memory:  m.Status.Allocatable.Memory().Value(),
			storage: m.Status.Allocatable.Storage().Value(),
			pods:    m.Status.Allocatable.Pods().Value(),
		}

		self.Nodes[m.Name] = &dashboardNode{
			Address: m.Status.Addresses[0].Address,
			Status:  findNodeStatus(m),
			Roles:   findNodeRoles(m),
			Cpu: &dashboardNodeResource{
				Allocatable: humanize.Comma(r.cpu),
			},
			Memory: &dashboardNodeResource{
				Allocatable: humanize.Comma(r.memory / (1024 * 1024)),
			},
			Storage: &dashboardNodeResource{
				Allocatable: humanize.Comma(r.storage / (1024 * 1024)),
			},
			Pods: &dashboardNodeResource{
				Allocatable: humanize.Comma(r.pods),
			},
		}
		allocate[m.Name] = r
		allocateTotal.cpu = +r.cpu
		allocateTotal.memory += r.memory
		allocateTotal.storage += r.storage
		allocateTotal.pods += r.pods
	}

	//beta1 <NODE> metrics list
	versionedNodeMetrics, err := metricsClient.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{})
	fmt.Println(versionedNodeMetrics)
	if err != nil {
		return err
	}

	//versioned (beta1) <NODE>  metrics 변환
	nodeMetrics := &metricsapi.NodeMetricsList{}
	err = metricsV1beta1api.Convert_v1beta1_NodeMetricsList_To_metrics_NodeMetricsList(versionedNodeMetrics, nodeMetrics, nil)
	if err != nil {
		return err
	}

	// 리소스 Usage 입력,  Percent 계산
	usageTotal := resource{}
	usage := coreV1.ResourceList{}
	for _, m := range nodeMetrics.Items {
		fmt.Println(m.CreationTimestamp)
		m.Usage.DeepCopyInto(&usage)
		nd := self.Nodes[m.Name]
		r := allocate[m.Name]
		nd.Cpu.Usage = humanize.Comma(m.Usage.Cpu().MilliValue())
		nd.Cpu.Percent = percent(m.Usage.Cpu().MilliValue(), r.cpu, 2)
		nd.Memory.Usage = humanize.Comma(m.Usage.Memory().Value() / (1024 * 1024))
		nd.Memory.Percent = percent(m.Usage.Memory().Value(), r.memory, 2)
		nd.Storage.Usage = humanize.Comma(m.Usage.Storage().Value() / (1024 * 1024))
		nd.Storage.Percent = percent(m.Usage.Storage().Value(), r.storage, 2)

		usageTotal.cpu += m.Usage.Cpu().MilliValue()
		usageTotal.memory += m.Usage.Memory().Value()
		usageTotal.storage += m.Usage.Storage().Value()

	}

	//노드별 파드 수 & running 파드 수
	usagePods := map[string]int64{}
	ready = 0
	for _, m := range podList.Items {
		if m.Spec.NodeName != "" {
			usagePods[m.Spec.NodeName] = usagePods[m.Spec.NodeName] + 1
			if m.Status.Phase == coreV1.PodRunning {
				ready += 1
			}
		}
	}
	self.Workloads.Pod = dashboardAvailable{Available: len(podList.Items), Ready: ready}

	for n, p := range usagePods {
		self.Nodes[n].Pods.Usage = humanize.Comma(p)
		self.Nodes[n].Pods.Percent = percent(p, allocate[n].pods, 2)
		usageTotal.pods += p
	}

	self.Cluster = dashboardCluster{
		Nodes:   fmt.Sprintf("%d", len(nodeList.Items)),
		Cpu:     percent(usageTotal.cpu, allocateTotal.cpu, 0),
		Memory:  percent(usageTotal.memory, allocateTotal.memory, 0),
		Storage: percent(usageTotal.storage, allocateTotal.storage, 0),
		Pods:    percent(usageTotal.pods, allocateTotal.pods, 0),
	}

	return nil

}

func percent(a int64, b int64, decimal int) string {
	return fmt.Sprintf(fmt.Sprintf("%%.%df", decimal), lang.DivideRound(a, b, (decimal+2))*100)
}

const (
	LabelNodeRolePrefix = "node-role.kubernetes.io/"
	NodeLabelRole       = "kubernetes.io/role"
)

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
