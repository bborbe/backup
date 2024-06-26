// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"

	v1 "github.com/bborbe/backup/k8s/apis/backup.benjamin-borbe.de/v1"
)

// TargetLister helps list Targets.
// All objects returned here must be treated as read-only.
type TargetLister interface {
	// List lists all Targets in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.Target, err error)
	// Targets returns an object that can list and get Targets.
	Targets(namespace string) TargetNamespaceLister
	TargetListerExpansion
}

// targetLister implements the TargetLister interface.
type targetLister struct {
	indexer cache.Indexer
}

// NewTargetLister returns a new TargetLister.
func NewTargetLister(indexer cache.Indexer) TargetLister {
	return &targetLister{indexer: indexer}
}

// List lists all Targets in the indexer.
func (s *targetLister) List(selector labels.Selector) (ret []*v1.Target, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Target))
	})
	return ret, err
}

// Targets returns an object that can list and get Targets.
func (s *targetLister) Targets(namespace string) TargetNamespaceLister {
	return targetNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// TargetNamespaceLister helps list and get Targets.
// All objects returned here must be treated as read-only.
type TargetNamespaceLister interface {
	// List lists all Targets in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.Target, err error)
	// Get retrieves the Target from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.Target, error)
	TargetNamespaceListerExpansion
}

// targetNamespaceLister implements the TargetNamespaceLister
// interface.
type targetNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Targets in the indexer for a given namespace.
func (s targetNamespaceLister) List(selector labels.Selector) (ret []*v1.Target, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Target))
	})
	return ret, err
}

// Get retrieves the Target from the indexer for a given namespace and name.
func (s targetNamespaceLister) Get(name string) (*v1.Target, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("target"), name)
	}
	return obj.(*v1.Target), nil
}
