package model

// base code : https://github.com/kubernetes/dashboard/tree/master/src/app/backend/resource/daemonset

import (
	"context"
	"fmt"

	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// returns a subset of pods controlled by given deployment.
func GetNodePods(apiClient *kubernetes.Clientset, name string) ([]v1.Pod, error) {

	podList, err := apiClient.CoreV1().Pods("").List(context.TODO(), metaV1.ListOptions{FieldSelector: fmt.Sprintf("spec.nodeName=%s", name)})
	if err != nil {
		return nil, err
	}
	return podList.Items, nil

}
