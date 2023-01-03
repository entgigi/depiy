package clients

import (
	"github.com/entgigi/depiy/operators/bundle-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type BundleV1Alpha1Api interface {
	Bundles(ns string) BundleInterface
}

type BundleV1Alpha1Client struct {
	client rest.Interface
}

func NewForConfig(c *rest.Config) (*BundleV1Alpha1Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1alpha1.GroupVersion.Group, Version: v1alpha1.GroupVersion.Version}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &BundleV1Alpha1Client{client}, nil
}

func (api *BundleV1Alpha1Client) Bundles(ns string) BundleInterface {
	return &bundleClient{
		restClient: api.client,
		namespace:  ns,
	}
}
