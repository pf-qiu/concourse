package present

import (
	"github.com/pf-qiu/concourse/v6/atc"
	"github.com/pf-qiu/concourse/v6/atc/db"
)

func VersionedResourceTypes(savedResourceTypes db.ResourceTypes) atc.VersionedResourceTypes {
	return savedResourceTypes.Deserialize()
}
