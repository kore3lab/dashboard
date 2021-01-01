/**
"k8s.io/client-go/dynamic"
    관련소스 : https://github.com/kubernetes/client-go/tree/master/dynamic


Kubernetes API Concepts
    https://kubernetes.io/docs/reference/using-api/api-concepts/

활용예제 참조
    https://github.com/kubernetes/dashboard
        /src/app/backend/resource/deployment/deploy.go
        DeployAppFromFile() 함수
*/
package client

import (
	"context"
	"fmt"
	"io"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	restclient "k8s.io/client-go/rest"
)

// type resourceVerber struct {
type DynamicClient struct {
	config *restclient.Config
}

// RestfulClient 리턴
func NewDynamicClient(config *restclient.Config) DynamicClient {
	return DynamicClient{config: config}
}

// List
func (self *DynamicClient) List(namespaceSet bool, namespace string, group string, version string, kind string) (r *unstructured.UnstructuredList, err error) {

	// 실행
	dynamicClient, err := dynamic.NewForConfig(self.config)
	if err != nil {
		return
	}

	// 예:  schema.GroupVersionResource{Group: "networking.istio.io", Version: "v1alpha3", Resource: "virtualservices"}
	groupVersionResource := schema.GroupVersionResource{Group: group, Version: version, Resource: kind}

	if namespaceSet {
		r, err = dynamicClient.Resource(groupVersionResource).Namespace(namespace).List(context.TODO(), metaV1.ListOptions{})

	} else {
		r, err = dynamicClient.Resource(groupVersionResource).List(context.TODO(), metaV1.ListOptions{})
	}

	return r, err

}

// GET
func (self *DynamicClient) GET(namespaceSet bool, namespace string, group string, version string, kind string, name string) (r *unstructured.Unstructured, err error) {

	// 실행
	dynamicClient, err := dynamic.NewForConfig(self.config)
	if err != nil {
		return
	}

	// 예:  schema.GroupVersionResource{Group: "networking.istio.io", Version: "v1alpha3", Resource: "virtualservices"}
	groupVersionResource := schema.GroupVersionResource{Group: group, Version: version, Resource: kind}

	if namespaceSet {
		r, err = dynamicClient.Resource(groupVersionResource).Namespace(namespace).Get(context.TODO(), name, metaV1.GetOptions{})

	} else {
		r, err = dynamicClient.Resource(groupVersionResource).Get(context.TODO(), name, metaV1.GetOptions{})
	}

	return r, err

}

// Watch
func (self *DynamicClient) Watch(namespaceSet bool, namespace string, group string, version string, kind string, resourceVersion string) (r watch.Interface, err error) {

	// 실행
	dynamicClient, err := dynamic.NewForConfig(self.config)
	if err != nil {
		return
	}

	// 예:  schema.GroupVersionResource{Group: "networking.istio.io", Version: "v1alpha3", Resource: "virtualservices"}
	groupVersionResource := schema.GroupVersionResource{Group: group, Version: version, Resource: kind}
	opts := metaV1.ListOptions{Watch: true}
	if resourceVersion != "" {
		opts.ResourceVersion = resourceVersion
	}

	if namespaceSet {
		r, err = dynamicClient.Resource(groupVersionResource).Namespace(namespace).Watch(context.TODO(), opts)

	} else {
		r, err = dynamicClient.Resource(groupVersionResource).Watch(context.TODO(), opts)
	}

	return r, err

}

// DELETE
func (self *DynamicClient) DELETE(namespaceSet bool, namespace string, group string, version string, kind string, name string) (err error) {

	// 실행
	dynamicClient, err := dynamic.NewForConfig(self.config)
	if err != nil {
		return
	}

	// 예:  schema.GroupVersionResource{Group: "networking.istio.io", Version: "v1alpha3", Resource: "virtualservices"}
	groupVersionResource := schema.GroupVersionResource{Group: group, Version: version, Resource: kind}

	if namespaceSet {
		err = dynamicClient.Resource(groupVersionResource).Namespace(namespace).Delete(context.TODO(), name, metaV1.DeleteOptions{})
	} else {
		err = dynamicClient.Resource(groupVersionResource).Delete(context.TODO(), name, metaV1.DeleteOptions{})
	}

	return err

}

// POST
func (self *DynamicClient) POST(payload io.Reader, isUpdate bool) (output *unstructured.Unstructured, err error) {

	d := yaml.NewYAMLOrJSONDecoder(payload, 4096)
	for {
		// payload 읽기
		data := &unstructured.Unstructured{}
		if err = d.Decode(data); err != nil {
			if err == io.EOF {
				return output, err
			}
			return output, err
		}

		// version kind
		version := data.GetAPIVersion()
		kind := data.GetKind()

		gv, err := schema.ParseGroupVersion(version)
		if err != nil {
			gv = schema.GroupVersion{Version: version}
		}

		// api version 에 해당하는 resource 정보 조회
		discoveryClient, err := discovery.NewDiscoveryClientForConfig(self.config)
		if err != nil {
			return output, err
		}

		apiResourceList, err := discoveryClient.ServerResourcesForGroupVersion(version)
		if err != nil {
			return output, err
		}
		apiResources := apiResourceList.APIResources

		var resource *metaV1.APIResource
		for _, apiResource := range apiResources {
			if apiResource.Kind == kind && !strings.Contains(apiResource.Name, "/") {
				resource = &apiResource
				break
			}
		}
		if resource == nil {
			err = fmt.Errorf("unknown resource kind: %s", kind)
			return output, err
		}

		// 실행
		dynamicClient, err := dynamic.NewForConfig(self.config)
		if err != nil {
			return output, err
		}

		// 예:  schema.GroupVersionResource{Group: "networking.istio.io", Version: "v1alpha3", Resource: "virtualservices"}
		groupVersionResource := schema.GroupVersionResource{Group: gv.Group, Version: gv.Version, Resource: resource.Name}

		// update 인 경우 resourceVersion 을 조회 & 수정
		if isUpdate {
			r, err := dynamicClient.Resource(groupVersionResource).Namespace(data.GetNamespace()).Get(context.TODO(), data.GetName(), metaV1.GetOptions{})
			if err != nil {
				return output, err
			}
			data.SetResourceVersion(r.GetResourceVersion())
			if resource.Namespaced {
				output, err = dynamicClient.Resource(groupVersionResource).Namespace(data.GetNamespace()).Update(context.TODO(), data, metaV1.UpdateOptions{})
			} else {
				output, err = dynamicClient.Resource(groupVersionResource).Update(context.TODO(), data, metaV1.UpdateOptions{})
			}
		} else {
			if resource.Namespaced {
				output, err = dynamicClient.Resource(groupVersionResource).Namespace(data.GetNamespace()).Create(context.TODO(), data, metaV1.CreateOptions{})
			} else {
				output, err = dynamicClient.Resource(groupVersionResource).Create(context.TODO(), data, metaV1.CreateOptions{})
			}
		}

		return output, err

	}

}
