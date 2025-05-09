// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	context "context"
	time "time"

	apisbackupbenjaminborbedev1 "github.com/bborbe/backup/k8s/apis/backup.benjamin-borbe.de/v1"
	versioned "github.com/bborbe/backup/k8s/client/clientset/versioned"
	internalinterfaces "github.com/bborbe/backup/k8s/client/informers/externalversions/internalinterfaces"
	backupbenjaminborbedev1 "github.com/bborbe/backup/k8s/client/listers/backup.benjamin-borbe.de/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// TargetInformer provides access to a shared informer and lister for
// Targets.
type TargetInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() backupbenjaminborbedev1.TargetLister
}

type targetInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewTargetInformer constructs a new informer for Target type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewTargetInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredTargetInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredTargetInformer constructs a new informer for Target type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredTargetInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.BackupV1().Targets(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.BackupV1().Targets(namespace).Watch(context.TODO(), options)
			},
		},
		&apisbackupbenjaminborbedev1.Target{},
		resyncPeriod,
		indexers,
	)
}

func (f *targetInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredTargetInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *targetInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&apisbackupbenjaminborbedev1.Target{}, f.defaultInformer)
}

func (f *targetInformer) Lister() backupbenjaminborbedev1.TargetLister {
	return backupbenjaminborbedev1.NewTargetLister(f.Informer().GetIndexer())
}
