// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"

	v1 "github.com/bborbe/backup/k8s/client/clientset/versioned/typed/backup.benjamin-borbe.de/v1"
)

type FakeBackupV1 struct {
	*testing.Fake
}

func (c *FakeBackupV1) Targets(namespace string) v1.TargetInterface {
	return newFakeTargets(c, namespace)
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeBackupV1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
