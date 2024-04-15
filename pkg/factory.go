package pkg

import (
	"context"

	"github.com/bborbe/errors"
	"github.com/bborbe/run"
)

func CreateSetupResourceDefinition(
	kubeConfig string,
	trigger run.Fire,
) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		k8sConnector := NewK8sConnector(kubeConfig)
		if err := k8sConnector.SetupCustomResourceDefinition(ctx); err != nil {
			return errors.Wrap(ctx, err, "setup resource definition failed")
		}
		trigger.Fire()
		<-ctx.Done()
		return nil
	}
}
