package model

import (
	"context"

	"github.com/kore3lab/dashboard/pkg/lang"
	"k8s.io/api/core/v1"
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

func ToPodList(pods []v1.Pod, metricsClient *versioned.Clientset) []Pod {

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
func GetPodsMatchLabels(k8sClient *kubernetes.Clientset, namespace string, selector labels.Selector) (*v1.PodList, error) {

	podList, err := k8sClient.CoreV1().Pods(namespace).List(context.TODO(), metaV1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		return nil, err
	}

	return podList, nil

}
