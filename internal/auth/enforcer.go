package auth

import (
	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

func NewEnforcer(db *gorm.DB) (*casbin.Enforcer, error) {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer("internal/conf/casbin.conf", adapter)
	if err != nil {
		return nil, err
	}

	err = enforcer.LoadPolicy()
	if err != nil {
		return nil, err
	}

	return enforcer, nil
}
