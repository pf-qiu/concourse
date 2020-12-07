package exec

import (
	"context"

	"github.com/pf-qiu/concourse/v6/atc/util"
	"github.com/hashicorp/go-multierror"
)

// AggregateStep is a step of steps to run in parallel.
type AggregateStep []Step

// Run executes all steps in parallel. It will indicate that it's ready when
// all of its steps are ready, and propagate any signal received to all running
// steps.
//
// It will wait for all steps to exit, even if one step fails or errors. After
// all steps finish, their errors (if any) will be aggregated and returned as a
// single error.
func (step AggregateStep) Run(ctx context.Context, state RunState) (bool, error) {
	oks := make(chan bool, len(step))
	errs := make(chan error, len(step))

	for _, s := range step {
		s := s
		go func() {
			defer func() {
				err := util.DumpPanic(recover(), "aggregate step")
				if err != nil {
					errs <- err
				}
			}()

			ok, err := s.Run(ctx, state)
			oks <- ok
			errs <- err
		}()
	}

	var result error
	for i := 0; i < len(step); i++ {
		err := <-errs
		if err != nil {
			result = multierror.Append(result, err)
		}
	}

	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	if result != nil {
		return false, result
	}

	ok := true
	for i := 0; i < len(step); i++ {
		ok = ok && <-oks
	}

	return ok, nil
}
