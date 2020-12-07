package secretsmanager

import (
	"github.com/pf-qiu/concourse/v6/atc/creds"
	flags "github.com/jessevdk/go-flags"
)

type managerFactory struct{}

func init() {
	creds.Register("secretsmanager", NewManagerFactory())
}

func NewManagerFactory() creds.ManagerFactory {
	return &managerFactory{}
}
func (manager managerFactory) Health() (interface{}, error) {
	return nil, nil
}

func (factory *managerFactory) AddConfig(group *flags.Group) creds.Manager {
	manager := &Manager{}
	subGroup, err := group.AddGroup("AWS SecretsManager Credential Management", "", manager)
	if err != nil {
		panic(err)
	}
	subGroup.Namespace = "aws-secretsmanager"
	return manager
}

func (factory *managerFactory) NewInstance(interface{}) (creds.Manager, error) {
	return &Manager{}, nil
}
