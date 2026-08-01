package app

import (
	"tms-be/internal/conf"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// registerCustomValidators dipanggil SEKALI pas startup (dari TmsAppBootstrap).
// Nambahin tag binding:"role" yang validasinya baca langsung dari
// conf.AllRoles — jadi dto.go gak perlu hardcode daftar role di
// binding:"oneof=...". Nambah role baru di conf/roles.go otomatis
// ke-cover di sini, gak perlu ubah apapun di module manapun.
func registerCustomValidators() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		panic("failed to register custom validators: gin validator engine is not *validator.Validate")
	}

	_ = v.RegisterValidation("role", func(fl validator.FieldLevel) bool {
		return conf.RoleType(fl.Field().String()).Valid()
	})
}
