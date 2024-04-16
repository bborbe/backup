// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package k8s

import (
	"context"

	"github.com/bborbe/errors"
	"github.com/golang/glog"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s_kubernetes "k8s.io/client-go/kubernetes"
)

//counterfeiter:generate -o mocks/k8s-ingress-deployer.go --fake-name K8sIngressDeployer . IngressDeployer
type IngressDeployer interface {
	Deploy(ctx context.Context, ingress v1.Ingress) error
	Undeploy(ctx context.Context, namespace Namespace, name string) error
}

func NewIngressDeployer(
	clientset k8s_kubernetes.Interface,
) IngressDeployer {
	return &ingressDeployer{
		clientset: clientset,
	}
}

type ingressDeployer struct {
	clientset k8s_kubernetes.Interface
}

func (s *ingressDeployer) Deploy(ctx context.Context, ingress v1.Ingress) error {
	_, err := s.clientset.NetworkingV1().Ingresses(ingress.Namespace).Get(ctx, ingress.Name, metav1.GetOptions{})
	if err != nil {
		_, err = s.clientset.NetworkingV1().Ingresses(ingress.Namespace).Create(ctx, &ingress, metav1.CreateOptions{})
		if err != nil {
			return errors.Wrap(ctx, err, "create ingress failed")
		}
		glog.V(3).Infof("ingress %s created successful", ingress.Name)
		return nil
	}
	_, err = s.clientset.NetworkingV1().Ingresses(ingress.Namespace).Update(ctx, &ingress, metav1.UpdateOptions{})
	if err != nil {
		return errors.Wrap(ctx, err, "update ingress failed")
	}
	glog.V(3).Infof("ingress %s updated successful", ingress.Name)
	return nil

}

func (s *ingressDeployer) Undeploy(ctx context.Context, namespace Namespace, name string) error {
	_, err := s.clientset.NetworkingV1().Ingresses(namespace.String()).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		glog.V(3).Infof("ingress '%s' not found => skip", name)
		return nil
	}
	if err := s.clientset.NetworkingV1().Ingresses(namespace.String()).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		return err
	}
	glog.V(3).Infof("delete %s completed", name)
	return nil
}
