package casbin

import (
	"github.com/casbin/casbin/v2"
	gormAdapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

type Enforcer = casbin.Enforcer

// NewEnforcer
func NewEnforcer(db *gorm.DB, modelPath string) (*casbin.Enforcer, error) {
	adapter, err := gormAdapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		return nil, err
	}

	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	return enforcer, nil
}
