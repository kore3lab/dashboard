package model

// base code : https://github.com/kubernetes/dashboard/tree/master/src/app/backend/resource/replicaset

import (
	"context"

	"github.com/kore3lab/dashboard/pkg/lang"

	appsV1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

// returns a subset of pods controlled by given deployment.
func GetReplicaSetPods(apiClient *kubernetes.Clientset, namespace string, name string) ([]v1.Pod, *v1.PodSpec, error) {

	replicaset, err := apiClient.AppsV1().ReplicaSets(namespace).Get(context.TODO(), name, metaV1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}

	labelSelector := labels.SelectorFromSet(replicaset.Spec.Selector.MatchLabels)

	podList, err := GetPodsMatchLabels(apiClient, namespace, labelSelector)
	if err != nil {
		return nil, nil, err
	}
	return lang.FilterPodsByControllerRef(replicaset, podList.Items), &replicaset.Spec.Template.Spec, nil

}

// return a subset of replicasets by given labelSelector
func GetReplicaSetMatchLabels(k8sClient *kubernetes.Clientset, namespace string, labelSelector labels.Selector) (*appsV1.ReplicaSetList, error) {

	rsList, err := k8sClient.AppsV1().ReplicaSets(namespace).List(context.TODO(), metaV1.ListOptions{LabelSelector: labelSelector.String()})
	if err != nil {
		return nil, err
	}

	return rsList, nil

}
