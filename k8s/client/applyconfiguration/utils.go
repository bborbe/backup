// Code generated by applyconfiguration-gen. DO NOT EDIT.

package applyconfiguration

import (
	schema "k8s.io/apimachinery/pkg/runtime/schema"

	v1 "github.com/bborbe/backup/k8s/apis/backup.benjamin-borbe.de/v1"
	backupbenjaminborbedev1 "github.com/bborbe/backup/k8s/client/applyconfiguration/backup.benjamin-borbe.de/v1"
)

// ForKind returns an apply configuration type for the given GroupVersionKind, or nil if no
// apply configuration type exists for the given GroupVersionKind.
func ForKind(kind schema.GroupVersionKind) interface{} {
	switch kind {
	// Group=backup, Version=v1
	case v1.SchemeGroupVersion.WithKind("BackupSpec"):
		return &backupbenjaminborbedev1.BackupSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("Target"):
		return &backupbenjaminborbedev1.TargetApplyConfiguration{}

	}
	return nil
}
