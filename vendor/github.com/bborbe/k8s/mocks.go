// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package k8s

import (
	"k8s.io/client-go/kubernetes"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

//counterfeiter:generate -o mocks/k8s-interface.go --fake-name K8sInterface . Interface
type Interface interface {
	kubernetes.Interface
}

//counterfeiter:generate -o mocks/k8s-appsv1-interface.go --fake-name K8sAppsV1Interface . AppsV1Interface
type AppsV1Interface interface {
	appsv1.AppsV1Interface
}

//counterfeiter:generate -o mocks/k8s-deployment-interface.go --fake-name K8sDeploymentInterface . DeploymentInterface
type DeploymentInterface interface {
	appsv1.DeploymentInterface
}

//counterfeiter:generate -o mocks/k8s-statefulset-interface.go --fake-name K8sStatefulSetInterface . StatefulSetInterface
type StatefulSetInterface interface {
	appsv1.StatefulSetInterface
}
