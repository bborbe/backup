// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"context"
	json "encoding/json"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"

	v1 "github.com/bborbe/backup/k8s/apis/backup.benjamin-borbe.de/v1"
	backupbenjaminborbedev1 "github.com/bborbe/backup/k8s/client/applyconfiguration/backup.benjamin-borbe.de/v1"
	scheme "github.com/bborbe/backup/k8s/client/clientset/versioned/scheme"
)

// TargetsGetter has a method to return a TargetInterface.
// A group's client should implement this interface.
type TargetsGetter interface {
	Targets(namespace string) TargetInterface
}

// TargetInterface has methods to work with Target resources.
type TargetInterface interface {
	Create(ctx context.Context, target *v1.Target, opts metav1.CreateOptions) (*v1.Target, error)
	Update(ctx context.Context, target *v1.Target, opts metav1.UpdateOptions) (*v1.Target, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Target, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.TargetList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Target, err error)
	Apply(ctx context.Context, target *backupbenjaminborbedev1.TargetApplyConfiguration, opts metav1.ApplyOptions) (result *v1.Target, err error)
	TargetExpansion
}

// targets implements TargetInterface
type targets struct {
	client rest.Interface
	ns     string
}

// newTargets returns a Targets
func newTargets(c *MonitoringV1Client, namespace string) *targets {
	return &targets{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the target, and returns the corresponding target object, and an error if there is any.
func (c *targets) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.Target, err error) {
	result = &v1.Target{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("targets").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Targets that match those selectors.
func (c *targets) List(ctx context.Context, opts metav1.ListOptions) (result *v1.TargetList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.TargetList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("targets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested targets.
func (c *targets) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("targets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a target and creates it.  Returns the server's representation of the target, and an error, if there is any.
func (c *targets) Create(ctx context.Context, target *v1.Target, opts metav1.CreateOptions) (result *v1.Target, err error) {
	result = &v1.Target{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("targets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(target).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a target and updates it. Returns the server's representation of the target, and an error, if there is any.
func (c *targets) Update(ctx context.Context, target *v1.Target, opts metav1.UpdateOptions) (result *v1.Target, err error) {
	result = &v1.Target{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("targets").
		Name(target.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(target).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the target and deletes it. Returns an error if one occurs.
func (c *targets) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("targets").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *targets) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("targets").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched target.
func (c *targets) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Target, err error) {
	result = &v1.Target{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("targets").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied target.
func (c *targets) Apply(ctx context.Context, target *backupbenjaminborbedev1.TargetApplyConfiguration, opts metav1.ApplyOptions) (result *v1.Target, err error) {
	if target == nil {
		return nil, fmt.Errorf("target provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(target)
	if err != nil {
		return nil, err
	}
	name := target.Name
	if name == nil {
		return nil, fmt.Errorf("target.Name must be provided to Apply")
	}
	result = &v1.Target{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("targets").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
