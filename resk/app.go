package resk

import (
	_ "github.com/memo012/red-packet/resk/api"
	"github.com/memo012/red-packet/resk/infra"
	"github.com/memo012/red-packet/resk/infra/base"
)

func init() {
	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})
	infra.Register(&infra.WebApiStarter{})
	infra.Register(&base.GinServerStarter{})
}
