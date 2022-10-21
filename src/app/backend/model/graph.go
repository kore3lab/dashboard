package model

import (
	"context"
	"fmt"

	"github.com/kore3lab/dashboard/pkg/config"
	"github.com/kore3lab/dashboard/pkg/lang"
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

// get group versions
func getGroupVersion(client *config.ClientSet) (coreVersion string, appsVersion string, networkVersion string, err error) {
	if discoveryClient, err := client.NewDiscoveryClient(); err == nil {
		if groups, err := discoveryClient.ServerGroups(); err == nil {
			for _, g := range groups.Groups {
				if g.Name == "" {
					coreVersion = g.PreferredVersion.GroupVersion
				} else if g.Name == "apps" {
					appsVersion = g.PreferredVersion.GroupVersion
				} else if g.Name == "network" {
					networkVersion = g.PreferredVersion.GroupVersion
				}
			}
		}
	}
	return
}
func newHierarchyNode(ty v1.TypeMeta, obj v1.ObjectMeta, owner string) HierarchyNode {
	h := HierarchyNode{
		UID:        string(obj.UID),
		APIVersion: ty.APIVersion, Kind: ty.Kind,
		Name: obj.Name, Namespace: obj.Namespace,
	}
	if len(owner) > 0 {
		h.Owner = owner
	} else if len(obj.OwnerReferences) > 0 {
		h.Owner = string(obj.OwnerReferences[0].UID)
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

	// get group versions
	var coreVersion string
	var appsVersion string
	if coreVersion, appsVersion, _, err = getGroupVersion(client); err != nil {
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
			hierarchy[deploy.Namespace] = append(hierarchy[deploy.Namespace], newHierarchyNode(v1.TypeMeta{APIVersion: appsVersion, Kind: "Deployment"}, deploy.ObjectMeta, ""))
		}
	}
	//deamonsets
	if deamonsets, err := api.AppsV1().DaemonSets(namespace).List(context.TODO(), v1.ListOptions{}); err != nil {
		return nil, err
	} else {
		for _, daemonset := range deamonsets.Items {
			hierarchy[daemonset.Namespace] = append(hierarchy[daemonset.Namespace], newHierarchyNode(v1.TypeMeta{APIVersion: appsVersion, Kind: "DaemonSet"}, daemonset.ObjectMeta, ""))
		}
	}
	//replicasets
	if replicasets, err := api.AppsV1().ReplicaSets(namespace).List(context.TODO(), v1.ListOptions{}); err != nil {
		return nil, err
	} else {
		for _, replicaset := range replicasets.Items {
			hierarchy[replicaset.Namespace] = append(hierarchy[replicaset.Namespace], newHierarchyNode(v1.TypeMeta{APIVersion: appsVersion, Kind: "ReplicaSet"}, replicaset.ObjectMeta, ""))
		}
	}
	//pods
	if pods, err := api.CoreV1().Pods(namespace).List(context.TODO(), v1.ListOptions{}); err != nil {
		return nil, err
	} else {
		for _, pod := range pods.Items {
			hierarchy[pod.Namespace] = append(hierarchy[pod.Namespace], newHierarchyNode(v1.TypeMeta{APIVersion: coreVersion, Kind: "Pod"}, pod.ObjectMeta, ""))
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

	// get group versions
	var coreVersion string
	var networkVersion string
	if coreVersion, _, networkVersion, err = getGroupVersion(client); err != nil {
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

	// service->pods
	for _, svc := range svcList.Items {

		hierarchy[svc.Namespace] = append(hierarchy[svc.Namespace], newHierarchyNode(v1.TypeMeta{APIVersion: coreVersion, Kind: "Service"}, svc.ObjectMeta, ""))

		// get the pods relations
		if len(svc.Spec.Selector) > 0 {
			selector, _ := v1.LabelSelectorAsSelector(&v1.LabelSelector{
				MatchLabels: svc.Spec.Selector,
			})
			for _, pod := range podList.Items {
				if pod.Namespace == svc.Namespace && selector.Matches(labels.Set(pod.Labels)) {
					n := newHierarchyNode(v1.TypeMeta{APIVersion: coreVersion, Kind: "Pod"}, pod.ObjectMeta, string(svc.UID))
					n.Line = pod.Status.PodIP
					hierarchy[pod.Namespace] = append(hierarchy[pod.Namespace], n)
				}
			}
		}

	}

	// ingress->services
	for _, ing := range ingList.Items {
		hierarchy[ing.Namespace] = append(hierarchy[ing.Namespace], newHierarchyNode(v1.TypeMeta{APIVersion: networkVersion, Kind: "Ingress"}, ing.ObjectMeta, ""))
		for _, rule := range ing.Spec.Rules {
			for _, path := range rule.HTTP.Paths {
				for i, nd := range hierarchy[ing.Namespace] {
					if nd.Kind == "Service" && path.Backend.Service.Name == nd.Name {
						if len(path.Path) > 1 {
							nd.Line = fmt.Sprintf("%s (%s) ", lang.NVL(rule.Host, "-"), path.Path)
						} else {
							nd.Line = rule.Host
						}
						nd.Owner = string(ing.UID)
						hierarchy[nd.Namespace] = append(hierarchy[nd.Namespace][:i], append(hierarchy[nd.Namespace][i+1:], nd)...)
					}
				}
			}
		}
	}

	return hierarchy, err
}

// workload graph
func GetPodGraph(cluster string, namespace string, name string) (Hierarchy, error) {

	// api-client
	client, err := config.Cluster.Client(cluster)
	if err != nil {
		return nil, err
	}

	// get group versions
	var coreVersion string
	var appsVersion string
	if coreVersion, appsVersion, _, err = getGroupVersion(client); err != nil {
		return nil, err
	}

	// get objects
	nodes := []HierarchyNode{}
	api, err := client.NewKubernetesClient()
	if err != nil {
		return nil, err
	}

	//pods
	var pod *coreV1.Pod
	if pod, err = api.CoreV1().Pods(namespace).Get(context.TODO(), name, v1.GetOptions{}); err != nil {
		return nil, err
	} else {
		n := newHierarchyNode(v1.TypeMeta{APIVersion: coreVersion, Kind: "Pod"}, pod.ObjectMeta, "")
		n.Owner = ""
		nodes = append(nodes, n)
	}

	//ownerReferences (ReplicaSet+Deployment, DaemonSet, StatefulSet)
	if len(pod.OwnerReferences) > 0 {
		owner := pod.OwnerReferences[0]

		if owner.Kind == "ReplicaSet" {
			if replicaset, err := api.AppsV1().ReplicaSets(namespace).Get(context.TODO(), owner.Name, v1.GetOptions{}); err != nil {
				return nil, err
			} else {
				n1 := newHierarchyNode(v1.TypeMeta{APIVersion: appsVersion, Kind: "ReplicaSet"}, replicaset.ObjectMeta, "")
				n1.Owner = string(pod.UID)
				nodes = append(nodes, n1)

				//deployment
				if len(replicaset.OwnerReferences) > 0 {
					if deployment, err := api.AppsV1().Deployments(namespace).Get(context.TODO(), replicaset.OwnerReferences[0].Name, v1.GetOptions{}); err == nil {
						nodes = append(nodes, newHierarchyNode(v1.TypeMeta{APIVersion: appsVersion, Kind: "Deployment"}, deployment.ObjectMeta, string(replicaset.UID)))
					}
				}
			}
		} else if owner.Kind == "DaemonSet" {
			if daemonset, err := api.AppsV1().DaemonSets(namespace).Get(context.TODO(), owner.Name, v1.GetOptions{}); err != nil {
				return nil, err
			} else {
				nodes = append(nodes, newHierarchyNode(v1.TypeMeta{APIVersion: appsVersion, Kind: "DaemonSet"}, daemonset.ObjectMeta, string(pod.UID)))
			}
		} else if owner.Kind == "StatefulSet" {
			if statefulset, err := api.AppsV1().StatefulSets(namespace).Get(context.TODO(), owner.Name, v1.GetOptions{}); err != nil {
				return nil, err
			} else {
				nodes = append(nodes, newHierarchyNode(v1.TypeMeta{APIVersion: appsVersion, Kind: "StatefulSet"}, statefulset.ObjectMeta, string(pod.UID)))
			}
		}
	}
	//ServiceAccount
	if pod.Spec.ServiceAccountName != "" {
		if serviceaccount, err := api.CoreV1().ServiceAccounts(namespace).Get(context.TODO(), pod.Spec.ServiceAccountName, v1.GetOptions{}); err != nil {
			return nil, err
		} else {
			nodes = append(nodes, newHierarchyNode(v1.TypeMeta{APIVersion: coreVersion, Kind: "ServiceAccount"}, serviceaccount.ObjectMeta, string(pod.UID)))
		}
	}

	//service
	if svcList, err := api.CoreV1().Services(namespace).List(context.TODO(), v1.ListOptions{}); err != nil {
		return nil, err
	} else {
		for _, svc := range svcList.Items {
			if len(svc.Spec.Selector) > 0 {
				selector, _ := v1.LabelSelectorAsSelector(&v1.LabelSelector{
					MatchLabels: svc.Spec.Selector,
				})
				if selector.Matches(labels.Set(pod.Labels)) {
					nodes = append(nodes, newHierarchyNode(v1.TypeMeta{APIVersion: coreVersion, Kind: "Service"}, svc.ObjectMeta, string(pod.UID)))
					break
				}
			}

		}
	}

	return map[string][]HierarchyNode{namespace: nodes}, err

}
