package model

// base code : https://github.com/kubernetes/dashboard/tree/master/src/app/backend/resource/daemonset

import (
	"context"

	"github.com/kore3lab/dashboard/pkg/lang"

	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// returns a subset of pods controlled by given deployment.
func GetDeploymentPods(apiClient *kubernetes.Clientset, namespace string, name string) ([]v1.Pod, *v1.PodSpec, error) {

	deployment, err := apiClient.AppsV1().Deployments(namespace).Get(context.TODO(), name, metaV1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}

	//https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#resources-that-support-set-based-requirements
	labelSelector, err := metaV1.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		return nil, nil, err
	}

	reList, err := GetReplicaSetMatchLabels(apiClient, deployment.GetNamespace(), labelSelector)
	if err != nil {
		return nil, nil, err
	}

	podList, err := GetPodsMatchLabels(apiClient, deployment.GetNamespace(), labelSelector)
	if err != nil {
		return nil, nil, err
	}

	return lang.FilterDeploymentPodsByOwnerReference(*deployment, reList.Items, podList.Items), &deployment.Spec.Template.Spec, nil

}
