package dummy

import (
	"path"
	"time"

	"github.com/pf-qiu/concourse/v6/atc/creds"
	"github.com/pf-qiu/concourse/v6/vars"
)

type Secrets struct {
	vars.StaticVariables

	TeamName     string
	PipelineName string
}

func (secrets *Secrets) NewSecretLookupPaths(teamName string, pipelineName string, allowRootPath bool) []creds.SecretLookupPath {
	lookupPaths := []creds.SecretLookupPath{}

	if len(pipelineName) > 0 {
		lookupPaths = append(lookupPaths, creds.NewSecretLookupWithPrefix(path.Join(teamName, pipelineName)+"/"))
	}

	lookupPaths = append(lookupPaths, creds.NewSecretLookupWithPrefix(teamName+"/"))
	lookupPaths = append(lookupPaths, creds.NewSecretLookupWithPrefix(""))

	return lookupPaths
}

func (secrets *Secrets) Get(secretPath string) (interface{}, *time.Time, bool, error) {
	v, found, err := secrets.StaticVariables.Get(vars.Reference{Path: secretPath})
	if err != nil {
		return nil, nil, false, err
	}

	if found {
		return v, nil, true, nil
	}

	return nil, nil, false, nil
}
