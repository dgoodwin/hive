// Code generated by main. DO NOT EDIT.

package v1

import (
	"time"

	v1 "github.com/openshift/hive/pkg/apis/hive/v1"
	scheme "github.com/openshift/hive/pkg/client/clientset-generated/clientset/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ClusterLeasePoolsGetter has a method to return a ClusterLeasePoolInterface.
// A group's client should implement this interface.
type ClusterLeasePoolsGetter interface {
	ClusterLeasePools() ClusterLeasePoolInterface
}

// ClusterLeasePoolInterface has methods to work with ClusterLeasePool resources.
type ClusterLeasePoolInterface interface {
	Create(*v1.ClusterLeasePool) (*v1.ClusterLeasePool, error)
	Update(*v1.ClusterLeasePool) (*v1.ClusterLeasePool, error)
	UpdateStatus(*v1.ClusterLeasePool) (*v1.ClusterLeasePool, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*v1.ClusterLeasePool, error)
	List(opts metav1.ListOptions) (*v1.ClusterLeasePoolList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.ClusterLeasePool, err error)
	ClusterLeasePoolExpansion
}

// clusterLeasePools implements ClusterLeasePoolInterface
type clusterLeasePools struct {
	client rest.Interface
}

// newClusterLeasePools returns a ClusterLeasePools
func newClusterLeasePools(c *HiveV1Client) *clusterLeasePools {
	return &clusterLeasePools{
		client: c.RESTClient(),
	}
}

// Get takes name of the clusterLeasePool, and returns the corresponding clusterLeasePool object, and an error if there is any.
func (c *clusterLeasePools) Get(name string, options metav1.GetOptions) (result *v1.ClusterLeasePool, err error) {
	result = &v1.ClusterLeasePool{}
	err = c.client.Get().
		Resource("clusterleasepools").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ClusterLeasePools that match those selectors.
func (c *clusterLeasePools) List(opts metav1.ListOptions) (result *v1.ClusterLeasePoolList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.ClusterLeasePoolList{}
	err = c.client.Get().
		Resource("clusterleasepools").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested clusterLeasePools.
func (c *clusterLeasePools) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("clusterleasepools").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a clusterLeasePool and creates it.  Returns the server's representation of the clusterLeasePool, and an error, if there is any.
func (c *clusterLeasePools) Create(clusterLeasePool *v1.ClusterLeasePool) (result *v1.ClusterLeasePool, err error) {
	result = &v1.ClusterLeasePool{}
	err = c.client.Post().
		Resource("clusterleasepools").
		Body(clusterLeasePool).
		Do().
		Into(result)
	return
}

// Update takes the representation of a clusterLeasePool and updates it. Returns the server's representation of the clusterLeasePool, and an error, if there is any.
func (c *clusterLeasePools) Update(clusterLeasePool *v1.ClusterLeasePool) (result *v1.ClusterLeasePool, err error) {
	result = &v1.ClusterLeasePool{}
	err = c.client.Put().
		Resource("clusterleasepools").
		Name(clusterLeasePool.Name).
		Body(clusterLeasePool).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *clusterLeasePools) UpdateStatus(clusterLeasePool *v1.ClusterLeasePool) (result *v1.ClusterLeasePool, err error) {
	result = &v1.ClusterLeasePool{}
	err = c.client.Put().
		Resource("clusterleasepools").
		Name(clusterLeasePool.Name).
		SubResource("status").
		Body(clusterLeasePool).
		Do().
		Into(result)
	return
}

// Delete takes name of the clusterLeasePool and deletes it. Returns an error if one occurs.
func (c *clusterLeasePools) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("clusterleasepools").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *clusterLeasePools) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("clusterleasepools").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched clusterLeasePool.
func (c *clusterLeasePools) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.ClusterLeasePool, err error) {
	result = &v1.ClusterLeasePool{}
	err = c.client.Patch(pt).
		Resource("clusterleasepools").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
