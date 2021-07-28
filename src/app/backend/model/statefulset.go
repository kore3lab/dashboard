package model

import (
	"context"

	"github.com/kore3lab/dashboard/pkg/lang"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

// returns a subset of pods controlled by given deployment.
func GetStatefulSetPods(apiClient *kubernetes.Clientset, namespace string, name string) ([]coreV1.Pod, *coreV1.PodSpec, error) {

	statefulset, err := apiClient.AppsV1().StatefulSets(namespace).Get(context.TODO(), name, metaV1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}

	labelSelector := labels.SelectorFromSet(statefulset.Spec.Selector.MatchLabels)

	podList, err := GetPodsMatchLabels(apiClient, namespace, labelSelector)
	if err != nil {
		return nil, nil, err
	}
	return lang.FilterPodsByControllerRef(statefulset, podList.Items), &statefulset.Spec.Template.Spec, nil

}

// statefulset's available-ready count in a cluster
func GetStatefulSetsReady(apiClient *kubernetes.Clientset, options metaV1.ListOptions) (available int, ready int, err error) {

	list, err := apiClient.AppsV1().StatefulSets("").List(context.TODO(), options)
	if err != nil {
		return available, ready, err
	}
	available = len(list.Items)
	for _, m := range list.Items {
		if m.Status.Replicas == m.Status.ReadyReplicas {
			ready += 1
		}
	}
	return available, ready, err

}
