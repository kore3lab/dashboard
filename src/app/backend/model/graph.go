package model

import (
	"context"
	"fmt"

	"github.com/kore3lab/dashboard/pkg/config"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetTopologyGraph(cluster string, namespace string) (topology Topology, err error) {

	topology = Topology{Nodes: []topologyNode{}, Links: []topologyLink{}}

	// api-client
	client, err := config.Cluster.Client(cluster)
	if err != nil {
		return
	}

	api, err := client.NewKubernetesClient()
	if err != nil {
		return
	}

	// pod
	podList, err := api.CoreV1().Pods(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return
	}
	for i := range podList.Items {
		pod := podList.Items[i]
		podID := fmt.Sprintf("%s", pod.ObjectMeta.UID)
		// append pod
		if pod.Spec.NodeName != "" {
			topology.Nodes = append(topology.Nodes,
				topologyNode{
					Id:        podID,
					Name:      pod.ObjectMeta.Name,
					Kind:      ELEMENT_KIND_POD,
					Namespace: pod.ObjectMeta.Namespace,
					Group:     pod.Spec.NodeName,
				})
			// linke  : pod --> node
			topology.Links = append(topology.Links,
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
				topology.Nodes = append(topology.Nodes,
					topologyNode{
						Id:        containerID,
						Name:      status.Image,
						Kind:      ELEMENT_KIND_CONTAINER,
						Namespace: pod.ObjectMeta.Namespace,
						Group:     pod.Spec.NodeName,
					})
				// linke  : container --> pod
				topology.Links = append(topology.Links,
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
	clsuterID := fmt.Sprintf("cluster-%s", cluster)
	topology.Nodes = append(topology.Nodes,
		topologyNode{
			Id:    clsuterID,
			Name:  cluster,
			Kind:  ELEMENT_KIND_CLUSTER,
			Group: "",
		})

	//nodes
	nodeList, err := api.CoreV1().Nodes().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return
	}
	for i := range nodeList.Items {
		node := nodeList.Items[i]
		nodeID := node.ObjectMeta.GetName()
		// append node
		topology.Nodes = append(topology.Nodes,
			topologyNode{
				Id:    nodeID,
				Name:  node.ObjectMeta.GetName(),
				Kind:  ELEMENT_KIND_NODE,
				Group: node.Name,
			})
		// linke  : node --> cluster
		topology.Links = append(topology.Links,
			topologyLink{
				Source: nodeID,
				Target: clsuterID,
				Hidden: false,
				Kind:   ELEMENT_KIND_NODE,
			})

	}

	return

}

func NewHierarchyNode(kind string, meta v1.ObjectMeta) HierarchyNode {
	h := HierarchyNode{
		HierarchyObject: HierarchyObject{Name: meta.Name, Namespace: meta.Namespace, Kind: kind},
	}
	if len(meta.OwnerReferences) > 0 {
		h.OwnerReference = HierarchyObject{Name: meta.OwnerReferences[0].Name, Namespace: meta.Namespace, Kind: meta.OwnerReferences[0].Kind}
	}

	return h
}

func GetWorkloadGraph(cluster string, namespace string) (Hierarchy, error) {

	hierarchy := make(map[string][]HierarchyNode)

	// api-client
	client, err := config.Cluster.Client(cluster)
	if err != nil {
		return nil, err
	}

	api, err := client.NewKubernetesClient()
	if err != nil {
		return nil, err
	}

	// namespace list
	var nsList *coreV1.NamespaceList
	if namespace == "" {
		if nsList, err = api.CoreV1().Namespaces().List(context.TODO(), v1.ListOptions{}); err != nil {
			return nil, err
		}
	} else {
		var ns *coreV1.Namespace
		if ns, err = api.CoreV1().Namespaces().Get(context.TODO(), namespace, v1.GetOptions{}); err != nil {
			return nil, err
		} else {
			nsList = &coreV1.NamespaceList{}
			nsList.Items = append(nsList.Items, *ns)
		}
	}

	for _, ns := range nsList.Items {

		namespace := ns.ObjectMeta.Name
		nodes := []HierarchyNode{}

		//deployment
		if deployments, err := api.AppsV1().Deployments(namespace).List(context.TODO(), v1.ListOptions{}); err != nil {
			return nil, err
		} else {
			for _, deploy := range deployments.Items {
				nodes = append(nodes, NewHierarchyNode("Deployment", deploy.ObjectMeta))
			}
		}
		//deamonsets
		if deamonsets, err := api.AppsV1().DaemonSets(namespace).List(context.TODO(), v1.ListOptions{}); err != nil {
			return nil, err
		} else {
			for _, daemonset := range deamonsets.Items {
				nodes = append(nodes, NewHierarchyNode("DaemonSet", daemonset.ObjectMeta))
			}
		}
		//replicasets
		if replicasets, err := api.AppsV1().ReplicaSets(namespace).List(context.TODO(), v1.ListOptions{}); err != nil {
			return nil, err
		} else {
			for _, replicaset := range replicasets.Items {
				nodes = append(nodes, NewHierarchyNode("ReplicaSet", replicaset.ObjectMeta))
			}
		}
		//pods
		if pods, err := api.CoreV1().Pods(namespace).List(context.TODO(), v1.ListOptions{}); err != nil {
			return nil, err
		} else {
			for _, pod := range pods.Items {
				nd := NewHierarchyNode("Pod", pod.ObjectMeta)
				nd.Depth = 2
				nodes = append(nodes, nd)
			}
		}
		hierarchy[namespace] = nodes
	}

	return hierarchy, err

}
