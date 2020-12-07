package noop

import "github.com/pf-qiu/concourse/v6/atc/creds"

type noopFactory struct{}

func NewNoopFactory() *noopFactory {
	return &noopFactory{}
}

func (*noopFactory) NewSecrets() creds.Secrets {
	return &Noop{}
}
