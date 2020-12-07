package concourse

import (
	"github.com/pf-qiu/concourse/v6/atc"
	"github.com/pf-qiu/concourse/v6/go-concourse/concourse/eventstream"
	"github.com/pf-qiu/concourse/v6/go-concourse/concourse/internal"
	"github.com/tedsuo/rata"
)

type Events interface {
	NextEvent() (atc.Event, error)
	Close() error
}

func (client *client) BuildEvents(buildID string) (Events, error) {
	sseEvents, err := client.connection.ConnectToEventStream(internal.Request{
		RequestName: atc.BuildEvents,
		Params: rata.Params{
			"build_id": buildID,
		},
	})
	if err != nil {
		return nil, err
	}

	return eventstream.NewSSEEventStream(sseEvents), nil
}
