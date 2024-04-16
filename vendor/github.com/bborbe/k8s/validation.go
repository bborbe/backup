// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package k8s

import (
	"context"

	"github.com/bborbe/errors"
	"github.com/bborbe/validation"
)

func NotEmptyString[T ~string](value T) validation.HasValidation {
	return validation.HasValidationFunc(func(ctx context.Context) error {
		if len(value) == 0 {
			return errors.Wrapf(ctx, validation.Error, "empty string")
		}
		return nil
	})
}

func NotEmptySlice[T any](values []T) validation.HasValidation {
	return validation.HasValidationFunc(func(ctx context.Context) error {
		if len(values) == 0 {
			return errors.Wrapf(ctx, validation.Error, "empty slice")
		}
		return nil
	})
}
