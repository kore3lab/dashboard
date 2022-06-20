package model

import (
	"context"
	"fmt"

	"github.com/kore3lab/dashboard/pkg/config"
	coreV1 "k8s.io/api/core/v1"
	networkV1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// topoplogy graph
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
		UID:  string(meta.UID),
		Name: meta.Name, Namespace: meta.Namespace, Kind: kind,
	}
	if len(meta.OwnerReferences) > 0 {
		h.Owner = string(meta.OwnerReferences[0].UID)
	}

	return h
}

// workload graph
func GetWorkloadGraph(cluster string, namespace string) (Hierarchy, error) {

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
	hierarchy := make(map[string][]HierarchyNode)
	if namespace == "" {
		if nsList, err := api.CoreV1().Namespaces().List(context.TODO(), v1.ListOptions{}); err != nil {
			return nil, err
		} else {
			for _, ns := range nsList.Items {
				hierarchy[ns.Name] = []HierarchyNode{}
			}
		}
	} else {
		hierarchy[namespace] = []HierarchyNode{}
	}

	//deployment
	if deployments, err := api.AppsV1().Deployments(namespace).List(context.TODO(), v1.ListOptions{}); err != nil {
		return nil, err
	} else {
		for _, deploy := range deployments.Items {
			hierarchy[deploy.Namespace] = append(hierarchy[deploy.Namespace], NewHierarchyNode("Deployment", deploy.ObjectMeta))
		}
	}
	//deamonsets
	if deamonsets, err := api.AppsV1().DaemonSets(namespace).List(context.TODO(), v1.ListOptions{}); err != nil {
		return nil, err
	} else {
		for _, daemonset := range deamonsets.Items {
			hierarchy[daemonset.Namespace] = append(hierarchy[daemonset.Namespace], NewHierarchyNode("DaemonSet", daemonset.ObjectMeta))
		}
	}
	//replicasets
	if replicasets, err := api.AppsV1().ReplicaSets(namespace).List(context.TODO(), v1.ListOptions{}); err != nil {
		return nil, err
	} else {
		for _, replicaset := range replicasets.Items {
			hierarchy[replicaset.Namespace] = append(hierarchy[replicaset.Namespace], NewHierarchyNode("ReplicaSet", replicaset.ObjectMeta))
		}
	}
	//pods
	if pods, err := api.CoreV1().Pods(namespace).List(context.TODO(), v1.ListOptions{}); err != nil {
		return nil, err
	} else {
		for _, pod := range pods.Items {
			hierarchy[pod.Namespace] = append(hierarchy[pod.Namespace], NewHierarchyNode("Pod", pod.ObjectMeta))
		}
	}

	return hierarchy, err

}

// network graph
func GetNetworkGraph(cluster string, namespace string) (Hierarchy, error) {

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
	hierarchy := make(map[string][]HierarchyNode)
	if namespace == "" {
		if nsList, err := api.CoreV1().Namespaces().List(context.TODO(), v1.ListOptions{}); err != nil {
			return nil, err
		} else {
			for _, ns := range nsList.Items {
				hierarchy[ns.Name] = []HierarchyNode{}
			}
		}
	} else {
		hierarchy[namespace] = []HierarchyNode{}
	}

	//service
	var svcList *coreV1.ServiceList
	if svcList, err = api.CoreV1().Services(namespace).List(context.TODO(), v1.ListOptions{}); err != nil {
		return nil, err
	}
	//pods
	var podList *coreV1.PodList
	if podList, err = api.CoreV1().Pods(namespace).List(context.TODO(), v1.ListOptions{}); err != nil {
		return nil, err
	}

	//ingress
	var ingList *networkV1.IngressList
	if ingList, err = api.NetworkingV1().Ingresses(namespace).List(context.TODO(), v1.ListOptions{}); err != nil {
		return nil, err
	}

	// service-pods
	for _, svc := range svcList.Items {
		selector, _ := v1.LabelSelectorAsSelector(&v1.LabelSelector{
			MatchLabels: svc.Spec.Selector,
		})
		hierarchy[svc.Namespace] = append(hierarchy[svc.Namespace], NewHierarchyNode("Service", svc.ObjectMeta))

		for _, pod := range podList.Items {
			if pod.Namespace == svc.Namespace && selector.Matches(labels.Set(pod.Labels)) {
				hierarchy[pod.Namespace] = append(hierarchy[pod.Namespace], HierarchyNode{
					UID: string(pod.UID), Name: pod.Name, Namespace: pod.Namespace, Kind: "Pod",
					Owner: string(svc.UID),
				})
			}
		}
	}

	// ingress-services
	for _, ing := range ingList.Items {
		for _, rule := range ing.Spec.Rules {
			for _, path := range rule.HTTP.Paths {
				for i, nd := range hierarchy[ing.Namespace] {
					if nd.Kind == "Service" && path.Backend.Service.Name == nd.Name && ing.Namespace == nd.Namespace {
						nd.Arrow = rule.Host
						nd.Owner = string(ing.UID)
						hierarchy[nd.Namespace] = append(hierarchy[nd.Namespace][:i], append(hierarchy[nd.Namespace][i+1:], nd)...)
					}
				}
			}
		}
	}

	return hierarchy, err

}
