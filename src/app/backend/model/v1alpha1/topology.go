package model

import (
	"context"
	"fmt"
	_ "time"

	"github.com/acornsoftlab/dashboard/pkg/config"
	"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/apimachinery/pkg/labels"
)

type Topology struct {
	context string
	Nodes   []topologyNode `json:"nodes"`
	Links   []topologyLink `json:"links"`
}
type topologyNode struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Kind   string `json:"kind"`
	Group  string `json:"group"`
	Labels string `json:"labels"`
}
type topologyLink struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Hidden bool   `json:"hidden"`
	Kind   string `json:"kind"`
}

const (
	ELEMENT_KIND_CLUSTER    string = "cluster"
	ELEMENT_KIND_NAMESPACE  string = "namespace"
	ELEMENT_KIND_POD        string = "pod"
	ELEMENT_KIND_NODE       string = "node"
	ELEMENT_KIND_REPLICASET string = "replica"
	ELEMENT_KIND_CONTAINER  string = "container"
)

func NewTopology(contextName string) Topology {
	return Topology{context: contextName, Nodes: []topologyNode{}}
}

func (self *Topology) Get(namespace string) error {

	// kubeconfig
	conf, err := config.KubeConfigs(self.context)
	if err != nil {
		return err
	}

	// api-client
	api, err := kubernetes.NewForConfig(conf)
	if err != nil {
		return err
	}

	// namespace := "cb-system"
	// namespace := ""

	// pod
	podList, err := api.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	for i := range podList.Items {
		pod := podList.Items[i]
		podID := fmt.Sprintf("%s", pod.ObjectMeta.UID)
		// append pod
		if pod.Spec.NodeName != "" {
			self.Nodes = append(self.Nodes,
				topologyNode{
					Id:    podID,
					Name:  pod.ObjectMeta.Name,
					Kind:  ELEMENT_KIND_POD,
					Group: pod.Spec.NodeName,
				})
			// linke  : pod --> node
			self.Links = append(self.Links,
				topologyLink{
					Source: podID,
					Target: pod.Spec.NodeName,
					Kind:   ELEMENT_KIND_POD,
					Hidden: false,
				})

			// append container
			for n := range pod.Status.ContainerStatuses {
				status := pod.Status.ContainerStatuses[n]
				containerID := status.ContainerID
				// append container
				self.Nodes = append(self.Nodes,
					topologyNode{
						Id:    containerID,
						Name:  status.Image,
						Kind:  ELEMENT_KIND_CONTAINER,
						Group: pod.Spec.NodeName,
					})
				// linke  : container --> pod
				self.Links = append(self.Links,
					topologyLink{
						Source: containerID,
						Target: podID,
						Kind:   ELEMENT_KIND_CONTAINER,
						Hidden: false,
					})
			}
		}

	}

	// cluster
	clsuterID := fmt.Sprintf("cluster-%s", self.context)
	self.Nodes = append(self.Nodes,
		topologyNode{
			Id:    clsuterID,
			Name:  self.context,
			Kind:  ELEMENT_KIND_CLUSTER,
			Group: "",
		})

	//nodes
	nodeList, err := api.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	for i := range nodeList.Items {
		node := nodeList.Items[i]
		nodeID := node.ObjectMeta.GetName()
		// append node
		self.Nodes = append(self.Nodes,
			topologyNode{
				Id:    nodeID,
				Name:  node.ObjectMeta.GetName(),
				Kind:  ELEMENT_KIND_NODE,
				Group: node.Name,
			})
		// linke  : node --> cluster
		self.Links = append(self.Links,
			topologyLink{
				Source: nodeID,
				Target: clsuterID,
				Hidden: false,
				Kind:   ELEMENT_KIND_NODE,
			})

	}

	return nil

}
