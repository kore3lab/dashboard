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
func GetJobPods(apiClient *kubernetes.Clientset, namespace string, name string) ([]v1.Pod, *v1.PodSpec, error) {

	job, err := apiClient.BatchV1().Jobs(namespace).Get(context.TODO(), name, metaV1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}

	labelSelector, err := metaV1.LabelSelectorAsSelector(job.Spec.Selector)
	if err != nil {
		return nil, nil, err
	}

	podList, err := GetPodsMatchLabels(apiClient, namespace, labelSelector)
	if err != nil {
		return nil, nil, err
	}
	return lang.FilterPodsByControllerRef(job, podList.Items), &job.Spec.Template.Spec, nil

}
