package lang

import (
	"fmt"

	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// returns a subset of pods controlled by given deployment.
// base code : https://github.com/kubernetes/dashboard/blob/master/src/app/backend/resource/common/pod.go
func FilterDeploymentPodsByOwnerReference(deployment appsV1.Deployment, allRS []appsV1.ReplicaSet, allPods []v1.Pod) []v1.Pod {
	var matchingPods []v1.Pod
	for _, rs := range allRS {
		if metaV1.IsControlledBy(&rs, &deployment) {
			matchingPods = append(matchingPods, FilterPodsByControllerRef(&rs, allPods)...)
		}
	}
	return matchingPods
}

// returns a subset of pods controlled by given controller resource, excluding deployments.
// base code : https://github.com/kubernetes/dashboard/blob/master/src/app/backend/resource/common/pod.go
func FilterPodsByControllerRef(owner metaV1.Object, allPods []v1.Pod) []v1.Pod {
	var matchingPods []v1.Pod
	for _, pod := range allPods {
		if metaV1.IsControlledBy(&pod, owner) {
			matchingPods = append(matchingPods, pod)
		}
	}
	return matchingPods
}

// get pod ready (2/2)
func GetPodReady(pod v1.Pod) string {

	ready := 0
	for _, v := range pod.Status.ContainerStatuses {
		if v.Ready == true {
			ready = ready + 1
		}
	}
	return fmt.Sprintf("%d/%d", ready, len(pod.Spec.Containers))
}

// getPodStatus returns status string calculated based on the same logic as kubectl
// Base code: https://github.com/kubernetes/kubernetes/blob/master/pkg/printers/internalversion/printers.go#L734
func GetPodStatus(pod v1.Pod) string {

	restarts := 0
	readyContainers := 0

	reason := string(pod.Status.Phase)
	if pod.Status.Reason != "" {
		reason = pod.Status.Reason
	}

	initializing := false
	for i := range pod.Status.InitContainerStatuses {
		container := pod.Status.InitContainerStatuses[i]
		restarts += int(container.RestartCount)
		switch {
		case container.State.Terminated != nil && container.State.Terminated.ExitCode == 0:
			continue
		case container.State.Terminated != nil:
			// initialization is failed
			if len(container.State.Terminated.Reason) == 0 {
				if container.State.Terminated.Signal != 0 {
					reason = fmt.Sprintf("Init: Signal %d", container.State.Terminated.Signal)
				} else {
					reason = fmt.Sprintf("Init: ExitCode %d", container.State.Terminated.ExitCode)
				}
			} else {
				reason = "Init:" + container.State.Terminated.Reason
			}
			initializing = true
		case container.State.Waiting != nil && len(container.State.Waiting.Reason) > 0 && container.State.Waiting.Reason != "PodInitializing":
			reason = fmt.Sprintf("Init: %s", container.State.Waiting.Reason)
			initializing = true
		default:
			reason = fmt.Sprintf("Init: %d/%d", i, len(pod.Spec.InitContainers))
			initializing = true
		}
		break
	}
	if !initializing {
		restarts = 0
		hasRunning := false
		for i := len(pod.Status.ContainerStatuses) - 1; i >= 0; i-- {
			container := pod.Status.ContainerStatuses[i]

			restarts += int(container.RestartCount)
			if container.State.Waiting != nil && container.State.Waiting.Reason != "" {
				reason = container.State.Waiting.Reason
			} else if container.State.Terminated != nil && container.State.Terminated.Reason != "" {
				reason = container.State.Terminated.Reason
			} else if container.State.Terminated != nil && container.State.Terminated.Reason == "" {
				if container.State.Terminated.Signal != 0 {
					reason = fmt.Sprintf("Signal: %d", container.State.Terminated.Signal)
				} else {
					reason = fmt.Sprintf("ExitCode: %d", container.State.Terminated.ExitCode)
				}
			} else if container.Ready && container.State.Running != nil {
				hasRunning = true
				readyContainers++
			}
		}

		// change pod status back to "Running" if there is at least one container still reporting as "Running" status
		if reason == "Completed" && hasRunning {
			if hasPodReadyCondition(pod.Status.Conditions) {
				reason = string(v1.PodRunning)
			} else {
				reason = "NotReady"
			}
		}
	}

	if pod.DeletionTimestamp != nil && pod.Status.Reason == "NodeLost" {
		reason = string(v1.PodUnknown)
	} else if pod.DeletionTimestamp != nil {
		reason = "Terminating"
	}

	if len(reason) == 0 {
		reason = string(v1.PodUnknown)
	}

	return reason

}

func hasPodReadyCondition(conditions []v1.PodCondition) bool {
	for _, condition := range conditions {
		if condition.Type == v1.PodReady && condition.Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}
