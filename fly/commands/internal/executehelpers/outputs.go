package executehelpers

import (
	"fmt"

	"github.com/pf-qiu/concourse/v6/atc"
	"github.com/pf-qiu/concourse/v6/fly/commands/internal/flaghelpers"
)

type Output struct {
	Name string
	Path string
	Plan atc.Plan
}

func DetermineOutputs(
	fact atc.PlanFactory,
	taskOutputs []atc.TaskOutputConfig,
	outputMappings []flaghelpers.OutputPairFlag,
) ([]Output, error) {
	outputs := []Output{}

	for _, i := range outputMappings {
		outputName := i.Name

		notInConfig := true
		for _, configOutput := range taskOutputs {
			if configOutput.Name == outputName {
				notInConfig = false
			}
		}

		if notInConfig {
			return nil, fmt.Errorf("unknown output '%s'", outputName)
		}

		outputs = append(outputs, Output{
			Name: outputName,
			Path: i.Path,
			Plan: fact.NewPlan(atc.ArtifactOutputPlan{
				Name: outputName,
			}),
		})
	}

	return outputs, nil
}
