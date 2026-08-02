package app

import (
	"tms-be/internal/conf"
	"tms-be/internal/sloc"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func registerCustomValidators() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		panic("failed to register custom validators: gin validator engine is not *validator.Validate")
	}

	_ = v.RegisterValidation("role", func(fl validator.FieldLevel) bool {
		return conf.RoleType(fl.Field().String()).Valid()
	})

	_ = v.RegisterValidation("slocMember", func(fl validator.FieldLevel) bool {
		return sloc.MemberType(fl.Field().String()).Valid()
	})

}
