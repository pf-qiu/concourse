package resource

import (
	"context"
	"fmt"
	"strings"

	"github.com/pf-qiu/concourse/v6/atc/runtime"
)

func (resource *resource) Put(
	ctx context.Context,
	spec runtime.ProcessSpec,
	runnable runtime.Runner,
) (runtime.VersionResult, error) {
	vr := runtime.VersionResult{}

	input, err := resource.Signature()
	if err != nil {
		return vr, err
	}

	err = runnable.RunScript(
		ctx,
		spec.Path,
		spec.Args,
		input,
		&vr,
		spec.StderrWriter,
		true,
	)
	if err != nil {
		return runtime.VersionResult{}, err
	}
	if vr.Version == nil {
		return runtime.VersionResult{}, fmt.Errorf("resource script (%s %s) output a null version", spec.Path, strings.Join(spec.Args, " "))
	}

	return vr, nil
}
