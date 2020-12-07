package commands

import (
	"fmt"

	"github.com/pf-qiu/concourse/v6/fly/commands/internal/flaghelpers"
	"github.com/pf-qiu/concourse/v6/fly/rc"
)

type LandWorkerCommand struct {
	Worker flaghelpers.WorkerFlag `short:"w"  long:"worker" required:"true" description:"Worker to land"`
}

func (command *LandWorkerCommand) Execute(args []string) error {
	workerName := command.Worker.Name()

	target, err := rc.LoadTarget(Fly.Target, Fly.Verbose)
	if err != nil {
		return err
	}

	err = target.Validate()
	if err != nil {
		return err
	}

	err = target.Client().LandWorker(workerName)
	if err != nil {
		return err
	}

	fmt.Printf("landed '%s'\n", workerName)

	return nil
}
