package resk

import (
	"github.com/memo012/red-packet/resk/infra"
	"github.com/memo012/red-packet/resk/infra/base"
)

func init()  {
	infra.Register(&base.PropsStarter{})
}