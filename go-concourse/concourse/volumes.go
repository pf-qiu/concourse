package concourse

import (
	"github.com/pf-qiu/concourse/v6/atc"
	"github.com/pf-qiu/concourse/v6/go-concourse/concourse/internal"
	"github.com/tedsuo/rata"
)

func (team *team) ListVolumes() ([]atc.Volume, error) {
	var volumes []atc.Volume

	params := rata.Params{
		"team_name": team.Name(),
	}
	err := team.connection.Send(internal.Request{
		RequestName: atc.ListVolumes,
		Params:      params,
	}, &internal.Response{
		Result: &volumes,
	})

	return volumes, err
}
