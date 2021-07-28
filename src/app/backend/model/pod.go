package model

import (
	"context"

	"github.com/kore3lab/dashboard/pkg/lang"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

//
type Pod struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Ready     string `json:"ready"`
	Status    string `json:"status"`
	Metrics   struct {
		CPU    int64 `json:"cpu"`
		Memory int64 `json:"memory"`
	} `json:"metrics"`
}

func ToPodList(pods []coreV1.Pod, metricsClient *versioned.Clientset) []Pod {

	podList := []Pod{}

	for _, pod := range pods {

		pm := Pod{
			Name:      pod.GetName(),
			Namespace: pod.GetNamespace(),
			Ready:     lang.GetPodReady(pod),
			Status:    lang.GetPodStatus(pod),
		}

		m, err := metricsClient.MetricsV1beta1().PodMetricses(pod.GetNamespace()).Get(context.TODO(), pod.Name, metaV1.GetOptions{})
		if err == nil {
			for _, c := range m.Containers {
				pm.Metrics.CPU = pm.Metrics.CPU + c.Usage.Cpu().MilliValue()
				pm.Metrics.Memory = pm.Metrics.Memory + c.Usage.Memory().Value()
			}
		}
		podList = append(podList, pm)
	}

	return podList
}

// return a subset of pods by given labelSelector
func GetPodsMatchLabels(k8sClient *kubernetes.Clientset, namespace string, selector labels.Selector) (*coreV1.PodList, error) {

	podList, err := k8sClient.CoreV1().Pods(namespace).List(context.TODO(), metaV1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		return nil, err
	}

	return podList, nil

}

// deployment's available-ready count in a cluster
func GetPodsReady(apiClient *kubernetes.Clientset, options metaV1.ListOptions) (available int, ready int, err error) {

	list, err := apiClient.CoreV1().Pods("").List(context.TODO(), options)
	if err != nil {
		return available, ready, err
	}
	available = len(list.Items)
	for _, m := range list.Items {
		if m.Spec.NodeName != "" {
			if m.Status.Phase == coreV1.PodRunning {
				ready += 1
			}
		}
	}

	return available, ready, err

}
