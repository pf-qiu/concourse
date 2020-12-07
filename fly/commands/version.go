package commands

import (
	"fmt"
	"os"

	"github.com/pf-qiu/concourse/v6/fly/rc"
)

func init() {
	Fly.Version = func() {
		fmt.Println(rc.LocalVersion)
		os.Exit(0)
	}
}
