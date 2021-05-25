/**
"k8s.io/client-go/dynamic"
    관련소스 : https://github.com/kubernetes/client-go/tree/master/dynamic


Kubernetes API Concepts
    https://kubernetes.io/docs/reference/using-api/api-concepts/

Patch
    https://kubernetes.io/docs/tasks/manage-kubernetes-objects/update-api-object-kubectl-patch/

활용예제 참조
    kubernetes-dashboard
        https://github.com/kubernetes/dashboard/blob/master/src/app/backend/resource/deployment/deploy.go
        DeployAppFromFile() 함수
*/
package client

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	// "encoding/json"
)

// type resourceVerber struct {
type DynamicClient struct {
	config       *rest.Config
	resource     schema.GroupVersionResource
	namespace    string
	namespaceSet bool
}

// RestfulClient 리턴
func NewDynamicClient(config *rest.Config) *DynamicClient {
	return &DynamicClient{
		config:       config,
		namespaceSet: false,
	}
}

// RestfulClient 리턴
func NewDynamicClientSchema(config *rest.Config, group string, version string, resource string) *DynamicClient {
	// 예:  schema.GroupVersionResource{Group: "networking.istio.io", Version: "v1alpha3", Resource: "virtualservices"}
	return &DynamicClient{
		config:       config,
		resource:     schema.GroupVersionResource{Group: group, Version: version, Resource: resource},
		namespaceSet: false,
	}
}

// List
func (self *DynamicClient) SetNamespace(namespace string) {
	self.namespace = namespace
	self.namespaceSet = (namespace != "")
}

// List
func (self *DynamicClient) List(opts v1.ListOptions) (r *unstructured.UnstructuredList, err error) {

	// 실행
	dynamicClient, err := dynamic.NewForConfig(self.config)
	if err != nil {
		return
	}

	if self.namespaceSet {
		r, err = dynamicClient.Resource(self.resource).Namespace(self.namespace).List(context.TODO(), opts)

	} else {
		r, err = dynamicClient.Resource(self.resource).List(context.TODO(), opts)
	}

	return r, err

}

// GET
func (self *DynamicClient) GET(name string, opts v1.GetOptions) (r *unstructured.Unstructured, err error) {

	// 실행
	dynamicClient, err := dynamic.NewForConfig(self.config)
	if err != nil {
		return
	}

	if self.namespaceSet {
		r, err = dynamicClient.Resource(self.resource).Namespace(self.namespace).Get(context.TODO(), name, opts)

	} else {
		r, err = dynamicClient.Resource(self.resource).Get(context.TODO(), name, opts)
	}

	return r, err

}

// Watch
func (self *DynamicClient) Watch(opts v1.ListOptions) (output watch.Interface, err error) {

	// 실행
	dynamicClient, err := dynamic.NewForConfig(self.config)
	if err != nil {
		return
	}
	opts.Watch = true

	if self.namespaceSet {
		output, err = dynamicClient.Resource(self.resource).Namespace(self.namespace).Watch(context.TODO(), opts)

	} else {
		output, err = dynamicClient.Resource(self.resource).Watch(context.TODO(), opts)
	}

	return output, err

}

// DELETE
func (self *DynamicClient) DELETE(name string, opts v1.DeleteOptions) (err error) {

	// 실행
	dynamicClient, err := dynamic.NewForConfig(self.config)
	if err != nil {
		return
	}

	if self.namespaceSet {
		err = dynamicClient.Resource(self.resource).Namespace(self.namespace).Delete(context.TODO(), name, opts)
	} else {
		err = dynamicClient.Resource(self.resource).Delete(context.TODO(), name, opts)
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

		var resource *v1.APIResource
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

		self.resource = schema.GroupVersionResource{Group: gv.Group, Version: gv.Version, Resource: resource.Name}
		self.namespace = data.GetNamespace()

		// update 인 경우 resourceVersion 을 조회 & 수정
		if isUpdate {
			r, err := dynamicClient.Resource(self.resource).Namespace(self.namespace).Get(context.TODO(), data.GetName(), v1.GetOptions{})
			if err != nil {
				return output, err
			}
			data.SetResourceVersion(r.GetResourceVersion())
			if resource.Namespaced {
				output, err = dynamicClient.Resource(self.resource).Namespace(self.namespace).Update(context.TODO(), data, v1.UpdateOptions{})
			} else {
				output, err = dynamicClient.Resource(self.resource).Update(context.TODO(), data, v1.UpdateOptions{})
			}
		} else {
			if resource.Namespaced {
				output, err = dynamicClient.Resource(self.resource).Namespace(self.namespace).Create(context.TODO(), data, v1.CreateOptions{})
			} else {
				output, err = dynamicClient.Resource(self.resource).Create(context.TODO(), data, v1.CreateOptions{})
			}
		}

		return output, err

	}

}

// Patch
func (self *DynamicClient) PATCH(name string, patchType types.PatchType, payload io.Reader, opts v1.PatchOptions) (output *unstructured.Unstructured, err error) {

	data, err := ioutil.ReadAll(payload)
	if err != nil {
		return output, err
	}

	// 실행
	dynamicClient, err := dynamic.NewForConfig(self.config)
	if err != nil {
		return output, err
	}

	if self.namespaceSet {
		output, err = dynamicClient.Resource(self.resource).Namespace(self.namespace).Patch(context.TODO(), name, patchType, data, opts)
	} else {
		output, err = dynamicClient.Resource(self.resource).Patch(context.TODO(), name, patchType, data, opts)
	}

	return output, err
}
