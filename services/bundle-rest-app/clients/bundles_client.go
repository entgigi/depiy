package clients

import (
	"context"

	"github.com/entgigi/depiy/operators/bundle-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

// copied from
// https://www.martin-helmich.de/en/blog/kubernetes-crd-client.html

type BundleInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.EntandoBundleV2List, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.EntandoBundleV2, error)
	//Create(*v1alpha1.EntandoBundleV2) (*v1alpha1.EntandoBundleV2, error)
}
type bundleClient struct {
	restClient rest.Interface
	namespace  string
}

func (bs *bundleClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.EntandoBundleV2, error) {
	result := v1alpha1.EntandoBundleV2{}
	err := bs.restClient.
		Get().
		Namespace(bs.namespace).
		Resource("entandobundlev2s").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (bs *bundleClient) List(opts metav1.ListOptions) (*v1alpha1.EntandoBundleV2List, error) {
	result := v1alpha1.EntandoBundleV2List{}
	err := bs.restClient.
		Get().
		Namespace(bs.namespace).
		Resource("entandobundlev2s").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}
