package model

import (
	"context"
	"fmt"
	"strings"

	resty "github.com/go-resty/resty/v2"
	"github.com/kore3lab/dashboard/pkg/config"
	"github.com/kore3lab/dashboard/pkg/lang"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	//log "github.com/sirupsen/logrus"
)

type ScraperClient struct {
	context   string
	apiClient *kubernetes.Clientset
	Limits    resources     `json:"limits"`
	Requests  resources     `json:"requests"`
	Metrics   []interface{} `json:"metrics"`
}

type resources struct {
	CPU    int64 `json:"cpu"`
	Memory int64 `json:"memory"`
}

func NewScraperClient(cluster string) (ScraperClient, error) {

	client, err := config.Cluster.Client(cluster)
	if err != nil {
		return ScraperClient{context: cluster}, err
	}

	apiClient, err := client.NewKubernetesClient()
	if err != nil {
		return ScraperClient{context: cluster}, err
	}

	return ScraperClient{context: cluster, apiClient: apiClient}, nil
}

// Get cluster metrics
func (self *ScraperClient) GetClusterMetrics() error {

	_, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetResult(&self.Metrics).
		Get(fmt.Sprintf("%s/api/v1/clusters/%s", config.Value.MetricsScraperUrl, self.context))

	return err

}

// Get node metrics
func (self *ScraperClient) GetNodeMetrics(node string) error {

	// invoke metrics-scraper api
	_, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetResult(&self.Metrics).
		Get(fmt.Sprintf("%s/api/v1/clusters/%s/nodes/%s", config.Value.MetricsScraperUrl, self.context, node))

	// get node allocatable resources
	nd, err := self.apiClient.CoreV1().Nodes().Get(context.TODO(), node, metaV1.GetOptions{})
	if err != nil {
		return err
	}

	self.Limits = resources{CPU: nd.Status.Allocatable.Cpu().MilliValue(), Memory: nd.Status.Allocatable.Memory().Value()}

	return err

}

// Get pod metrics
func (self *ScraperClient) GetPodMetrics(namespace string, pod string) error {

	// get requests, limits resources
	pd, err := self.apiClient.CoreV1().Pods(namespace).Get(context.TODO(), pod, metaV1.GetOptions{})
	if err != nil {
		return err
	}

	err = self.getPodMetrics(namespace, pd.GetObjectMeta().GetName(), nil, pd.Spec)
	if err != nil {
		return err
	}

	return nil

}

// Get deployment metrics
func (self *ScraperClient) GetDeploymentMetrics(namespace string, name string) error {

	// deployment
	deploy, err := self.apiClient.AppsV1().Deployments(namespace).Get(context.TODO(), name, metaV1.GetOptions{})
	if err != nil {
		return err
	}

	return self.getPodMetrics(namespace, "", deploy.Spec.Selector, deploy.Spec.Template.Spec)

}

// Get statefulset metrics
func (self *ScraperClient) GetSetatefulSetMetrics(namespace string, name string) error {

	// statefulset
	statefulset, err := self.apiClient.AppsV1().StatefulSets(namespace).Get(context.TODO(), name, metaV1.GetOptions{})
	if err != nil {
		return err
	}

	return self.getPodMetrics(namespace, "", statefulset.Spec.Selector, statefulset.Spec.Template.Spec)

}

// Get statefulset metrics
func (self *ScraperClient) GetDaemonSetMetrics(namespace string, name string) error {

	// statefulset
	daemonset, err := self.apiClient.AppsV1().DaemonSets(namespace).Get(context.TODO(), name, metaV1.GetOptions{})
	if err != nil {
		return err
	}

	return self.getPodMetrics(namespace, "", daemonset.Spec.Selector, daemonset.Spec.Template.Spec)

}

// Get replicaset metrics
func (self *ScraperClient) GetReplicaSetMetrics(namespace string, name string) error {

	// replicaset
	replicaset, err := self.apiClient.AppsV1().ReplicaSets(namespace).Get(context.TODO(), name, metaV1.GetOptions{})
	if err != nil {
		return err
	}

	return self.getPodMetrics(namespace, "", replicaset.Spec.Selector, replicaset.Spec.Template.Spec)

}

// Get pod metrics
func (self *ScraperClient) getPodMetrics(namespace string, name string, selector *metaV1.LabelSelector, podSpec coreV1.PodSpec) error {

	url := fmt.Sprintf("%s/api/v1/clusters/%s/namespaces/%s", config.Value.MetricsScraperUrl, self.context, namespace)

	if name == "" {
		// get child instanced pod list (if name is empty )
		podList, err := lang.GetPodsMatchLabels(self.apiClient, namespace, selector.MatchLabels)
		if err != nil {
			return err
		}
		names := []string{}
		for _, pd := range podList.Items {
			names = append(names, pd.ObjectMeta.Name)
		}
		url = url + fmt.Sprintf("/pods/%s/AVG", strings.Join(names, ","))
	} else {
		url = url + fmt.Sprintf("/pods/%s", name)
	}

	// invoke metrics-scraper api
	_, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetResult(&self.Metrics).
		Get(url)

	// get request, limit
	for _, c := range podSpec.Containers {
		self.Limits.CPU = self.Limits.CPU + c.Resources.Limits.Cpu().MilliValue()
		self.Limits.Memory = self.Limits.Memory + c.Resources.Limits.Memory().Value()
		self.Requests.CPU = self.Requests.CPU + c.Resources.Requests.Cpu().MilliValue()
		self.Requests.Memory = self.Requests.Memory + c.Resources.Requests.Memory().Value()
	}

	return err

}
