package casbin

import "github.com/casbin/casbin/v3"

type Service struct {
	enforcer *casbin.Enforcer
}

func New(enforcer *casbin.Enforcer) *Service {
	if enforcer == nil {
		panic("casbin enforcer is nil")
	}
	return &Service{enforcer: enforcer}
}

func (s *Service) Enforce(sub, obj, act string) (bool, error) {
	return s.enforcer.Enforce(sub, obj, act)
}

func (s *Service) GrantRole(email, role string) (bool, error) {
	return s.enforcer.AddGroupingPolicy(email, role)
}

func (s *Service) RevokeRole(email, role string) (bool, error) {
	return s.enforcer.RemoveGroupingPolicy(email, role)
}

func (s *Service) GrantPermission(role, obj, act string) (bool, error) {
	return s.enforcer.AddPolicy(role, obj, act)
}

func (s *Service) RevokePermission(role, obj, act string) (bool, error) {
	return s.enforcer.RemovePolicy(role, obj, act)
}

func (s *Service) GetRoleForUser(email string) ([]string, error) {
	return s.enforcer.GetRolesForUser(email)
}
