package worker

import (
	"context"

	"github.com/pf-qiu/concourse/v6/tsa"
)

//go:generate counterfeiter . TSAClient

type TSAClient interface {
	Register(context.Context, tsa.RegisterOptions) error

	Land(context.Context) error
	Retire(context.Context) error
	Delete(context.Context) error

	ReportContainers(context.Context, []string) error
	ContainersToDestroy(context.Context) ([]string, error)

	ReportVolumes(context.Context, []string) error
	VolumesToDestroy(context.Context) ([]string, error)
}
