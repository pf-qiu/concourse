package concourse

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/pf-qiu/concourse/v6/atc"
	"github.com/pf-qiu/concourse/v6/go-concourse/concourse/internal"
)

func (client *client) GetCLIReader(arch, platform string) (io.ReadCloser, http.Header, error) {
	responseHeaders := http.Header{}
	response := internal.Response{Headers: &responseHeaders}

	err := client.connection.Send(internal.Request{
		RequestName: atc.DownloadCLI,
		Query: url.Values{
			"arch":     {arch},
			"platform": {platform},
		},
		ReturnResponseBody: true,
	},
		&response,
	)
	if err != nil {
		return nil, nil, err
	}

	readCloser, ok := response.Result.(io.ReadCloser)
	if !ok {
		return nil, nil, errors.New("Unable to get stream from response.")
	}

	return readCloser, responseHeaders, nil
}
