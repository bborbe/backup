// Copyright (c) 2019 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg

import (
	"context"
	"time"

	"github.com/bborbe/errors"
	"github.com/bborbe/k8s"
	libk8s "github.com/bborbe/k8s"
	"github.com/golang/glog"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"

	backupv1 "github.com/bborbe/backup/k8s/apis/backup.benjamin-borbe.de/v1"
	"github.com/bborbe/backup/k8s/client/informers/externalversions"
)

const (
	defaultResync = 5 * time.Minute
	name          = "targets.backup.benjamin-borbe.de"
)

//counterfeiter:generate -o ../mocks/k8s-connector.go --fake-name K8sConnector . K8sConnector
type K8sConnector interface {
	SetupCustomResourceDefinition(ctx context.Context) error
	Listen(ctx context.Context, resourceEventHandler cache.ResourceEventHandler) error
	Targets(ctx context.Context) (backupv1.Targets, error)
	Target(ctx context.Context, name string) (*backupv1.Target, error)
}

func NewK8sConnector(
	backupClientset BackupClientset,
	apiextensionsInterface libk8s.ApiextensionsInterface,
	namespace k8s.Namespace,
) K8sConnector {
	return &k8sConnector{
		backupClientset:        backupClientset,
		apiextensionsInterface: apiextensionsInterface,
		namespace:              namespace,
	}
}

type k8sConnector struct {
	namespace              k8s.Namespace
	backupClientset        BackupClientset
	apiextensionsInterface libk8s.ApiextensionsInterface
}

func (k *k8sConnector) Target(ctx context.Context, name string) (*backupv1.Target, error) {
	target, err := k.backupClientset.BackupV1().Targets(k.namespace.String()).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(ctx, err, "list target failed")
	}
	return target, nil
}

func (k *k8sConnector) Targets(ctx context.Context) (backupv1.Targets, error) {
	targetList, err := k.backupClientset.BackupV1().Targets(k.namespace.String()).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(ctx, err, "list target failed")
	}
	return targetList.Items, nil
}

func (k *k8sConnector) Listen(
	ctx context.Context,
	resourceEventHandler cache.ResourceEventHandler,
) error {
	informerFactory := externalversions.NewSharedInformerFactory(k.backupClientset, defaultResync)
	_, err := informerFactory.
		Backup().
		V1().
		Targets().
		Informer().
		AddEventHandler(resourceEventHandler)
	if err != nil {
		return errors.Wrap(ctx, err, "add event handler failed")
	}

	stopCh := make(chan struct{})
	glog.V(2).Infof("listen for events")
	informerFactory.Start(stopCh)
	select {
	case <-ctx.Done():
		glog.V(0).Infof("listen canceled")
	case <-stopCh:
		glog.V(0).Infof("listen stopped")
	}
	return nil
}

func (k *k8sConnector) SetupCustomResourceDefinition(ctx context.Context) error {
	customResourceDefinition, err := k.apiextensionsInterface.ApiextensionsV1().CustomResourceDefinitions().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		glog.V(2).Infof("CustomResourceDefinition '%s' not found (%v) => create", name, err)
		if err := k.createCrd(ctx); err != nil {
			return errors.Wrap(ctx, err, "create crd failed")
		}
		return nil
	}
	if err := k.updateCrd(ctx, customResourceDefinition); err != nil {
		return errors.Wrap(ctx, err, "create crd failed")
	}
	return nil
}

func (k *k8sConnector) updateCrd(ctx context.Context, customResourceDefinition *v1.CustomResourceDefinition) error {
	customResourceDefinition.Spec = createSpec()
	if _, err := k.apiextensionsInterface.ApiextensionsV1().CustomResourceDefinitions().Update(ctx, customResourceDefinition, metav1.UpdateOptions{}); err != nil {
		return errors.Wrap(ctx, err, "update CustomResourceDefinition failed")
	}
	glog.V(2).Infof("CustomResourceDefinitions '%s' updated", name)
	return nil
}

func (k *k8sConnector) createCrd(ctx context.Context) error {
	_, err := k.apiextensionsInterface.ApiextensionsV1().CustomResourceDefinitions().Create(
		ctx,
		&v1.CustomResourceDefinition{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "apiextensions.k8s.io/v1",
				Kind:       "CustomResourceDefinition",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
			Spec: createSpec(),
		},
		metav1.CreateOptions{},
	)
	if err != nil {
		return errors.Wrap(ctx, err, "create CustomResourceDefinition failed")
	}
	glog.V(2).Infof("CustomResourceDefinition '%s' created", name)
	return nil
}

func boolPointer(value bool) *bool {
	return &value
}

func createSpec() v1.CustomResourceDefinitionSpec {
	return v1.CustomResourceDefinitionSpec{
		Group: "backup.benjamin-borbe.de",
		Names: v1.CustomResourceDefinitionNames{
			Kind:     "Target",
			ListKind: "TargetList",
			Plural:   "targets",
			Singular: "target",
		},
		Scope: "Namespaced",
		Versions: []v1.CustomResourceDefinitionVersion{
			{
				Name:    "v1",
				Served:  true,
				Storage: true,
				Schema: &v1.CustomResourceValidation{
					OpenAPIV3Schema: &v1.JSONSchemaProps{
						XPreserveUnknownFields: boolPointer(true),
					},
				},
			},
		},
	}
}
