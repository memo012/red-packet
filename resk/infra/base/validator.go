package base

import (
	"github.com/memo012/red-packet/resk/infra"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

func Validate() *validator.Validate {
	return validate
}

type ValidatorStarter struct {
	infra.BaseStarter
}

func (v *ValidatorStarter) Init(ctx infra.StarterContext) {
	validate = validator.New()
}
