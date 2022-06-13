package model

import (
	"context"
	"errors"
	"fmt"

	"github.com/kore3lab/dashboard/pkg/client"
	"github.com/kore3lab/dashboard/pkg/config"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// get cluster metrics
func GetClusterCumulativeMetrics(cluster string) (*NodeCumulativeMetrics, error) {

	clientSet, err := config.Cluster.Client(cluster)
	if err != nil {
		return nil, err
	}
	apiClient, err := clientSet.NewKubernetesClient()
	if err != nil {
		return nil, err
	}
	metricsClient := clientSet.NewCumulativeMetricsClient()

	result := NodeCumulativeMetrics{}

	// metrics
	metrics, err := metricsClient.Get(client.CumulativeMetricsResourceSelector{})
	if err != nil {
		return nil, err
	}
	result.Metrics = metrics

	// allocate, capacity
	nodeList, err := apiClient.CoreV1().Nodes().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, nd := range nodeList.Items {
		result.Allocatable.CPU += nd.Status.Allocatable.Cpu().MilliValue()
		result.Allocatable.Memory += nd.Status.Allocatable.Memory().Value()
		result.Capacity.CPU += nd.Status.Capacity.Cpu().MilliValue()
		result.Capacity.Memory += nd.Status.Capacity.Memory().Value()
	}

	return &result, nil

}

// get node metrics
func GetNodeCumulativeMetrics(cluster string, name string) (*NodeCumulativeMetrics, error) {

	clientSet, err := config.Cluster.Client(cluster)
	if err != nil {
		return nil, err
	}

	apiClient, err := clientSet.NewKubernetesClient()
	if err != nil {
		return nil, err
	}
	metricsClient := clientSet.NewCumulativeMetricsClient()

	node, err := apiClient.CoreV1().Nodes().Get(context.TODO(), name, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	result := NodeCumulativeMetrics{}
	metrics, err := metricsClient.Get(client.CumulativeMetricsResourceSelector{Node: name})
	if err != nil {
		return nil, err
	}
	result.Metrics = metrics

	result.Allocatable.CPU = node.Status.Allocatable.Cpu().MilliValue()
	result.Allocatable.Memory = node.Status.Allocatable.Memory().Value()
	result.Capacity.CPU = node.Status.Capacity.Cpu().MilliValue()
	result.Capacity.Memory = node.Status.Capacity.Memory().Value()

	return &result, nil
}

// get workloads metrics
func GetWorkloadCumulativeMetrics(cluster string, namespace string, resource string, name string) (*CumulativeMetrics, error) {

	clientSet, err := config.Cluster.Client(cluster)
	if err != nil {
		return nil, err
	}

	apiClient, err := clientSet.NewKubernetesClient()
	if err != nil {
		return nil, err
	}

	var pods []coreV1.Pod
	var podSpec *coreV1.PodSpec
	if resource == "pods" {
		pod, err := apiClient.CoreV1().Pods(namespace).Get(context.TODO(), name, metaV1.GetOptions{})
		if err != nil {
			return nil, err
		}
		pods = []coreV1.Pod{*pod}
		podSpec = &pod.Spec
	} else if resource == "deployments" {
		pods, podSpec, err = GetDeploymentPods(apiClient, namespace, name)
		if err != nil {
			return nil, err
		}
	} else if resource == "statefulsets" {
		pods, podSpec, err = GetStatefulSetPods(apiClient, namespace, name)
		if err != nil {
			return nil, err
		}
	} else if resource == "daemonsets" {
		pods, podSpec, err = GetDaemonSetPods(apiClient, namespace, name)
		if err != nil {
			return nil, err
		}
	} else if resource == "replicasets" {
		pods, podSpec, err = GetReplicaSetPods(apiClient, namespace, name)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New(fmt.Sprintf("unsupported resource '%s'", resource))
	}
	if len(pods) > 0 {
		names := []string{}
		for _, pd := range pods {
			names = append(names, pd.ObjectMeta.Name)
		}

		metricsClient := clientSet.NewCumulativeMetricsClient()

		selector := client.CumulativeMetricsResourceSelector{
			Pods:      names,
			Namespace: namespace,
			Function:  "AVG",
		}

		result := CumulativeMetrics{}
		// invoke metrics-scraper api
		metrics, err := metricsClient.Get(selector)
		if err != nil {
			return nil, err
		}
		result.Metrics = metrics

		// get request, limit
		for _, c := range podSpec.Containers {
			result.Limits.CPU = result.Limits.CPU + c.Resources.Limits.Cpu().MilliValue()
			result.Limits.Memory = result.Limits.Memory + c.Resources.Limits.Memory().Value()
			result.Requests.CPU = result.Requests.CPU + c.Resources.Requests.Cpu().MilliValue()
			result.Requests.Memory = result.Requests.Memory + c.Resources.Requests.Memory().Value()
		}

		return &result, nil
	} else {
		return &CumulativeMetrics{}, nil
	}

}

// get pod list with metrics
func GetNodePodListWithMetrics(cluster string, name string) (interface{}, error) {

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

	pods, err := GetNodePods(apiClient, name)
	if err != nil {
		return nil, err
	}

	return struct {
		Pods []Pod `json:"pods"`
	}{
		Pods: ToPodList(pods, metricsClient),
	}, nil

}

// get pod list with metrics
func GetWorkloadPodListWithMetrics(cluster string, namespace string, resource string, name string) (interface{}, error) {

	client, err := config.Cluster.Client(cluster)
	if err != nil {
		return nil, err
	}

	apiClient, err := client.NewKubernetesClient()
	if err != nil {
		return nil, err
	}

	metricsClient, err := client.NewMetricsClient()
	if err != nil {
		return nil, err
	}

	var pods []coreV1.Pod
	if resource == "deployments" {
		pods, _, err = GetDeploymentPods(apiClient, namespace, name)
		if err != nil {
			return nil, err
		}
	} else if resource == "statefulsets" {
		pods, _, err = GetStatefulSetPods(apiClient, namespace, name)
		if err != nil {
			return nil, err
		}
	} else if resource == "daemonsets" {
		pods, _, err = GetDaemonSetPods(apiClient, namespace, name)
		if err != nil {
			return nil, err
		}
	} else if resource == "replicasets" {
		pods, _, err = GetReplicaSetPods(apiClient, namespace, name)
		if err != nil {
			return nil, err
		}
	} else if resource == "jobs" {
		pods, _, err = GetJobPods(apiClient, namespace, name)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New(fmt.Sprintf("unsupported resource '%s'", resource))
	}

	return struct {
		Pods []Pod `json:"pods"`
	}{
		Pods: ToPodList(pods, metricsClient),
	}, nil

}
