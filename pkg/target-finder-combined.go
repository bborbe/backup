// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg

func NewCombinedTargetFinder(k8sConnector K8sConnector) TargetFinder {
	return TargetFinderList{
		NewTargetFinder(k8sConnector),           // Search by name first
		NewTargetFinderByHostname(k8sConnector), // Then search by hostname
	}
}
