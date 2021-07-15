package lang

import (
	"context"

	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

// Get pod list from matchlabels
func GetPodsMatchLabels(k8sClient *kubernetes.Clientset, namespace string, matchLabels map[string]string) (*v1.PodList, error) {

	labelSelector := labels.SelectorFromSet(matchLabels)
	podList, err := k8sClient.CoreV1().Pods(namespace).List(context.TODO(), metaV1.ListOptions{LabelSelector: labelSelector.String()})
	if err != nil {
		return nil, err
	}

	return podList, nil

}
