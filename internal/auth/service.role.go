package auth

import "github.com/casbin/casbin/v3"

type RoleService struct {
	enf *casbin.Enforcer
}

func NewRoleService(enf *casbin.Enforcer) *RoleService {
	if enf == nil {
		panic("auth.RoleService: nil enforcer")
	}
	return &RoleService{
		enf: enf,
	}
}

func (s *RoleService) Enforce(sub, obj, act string) (bool, error) {
	return s.enf.Enforce(sub, obj, act)
}

func (s *RoleService) GrantRole(email, role string) (bool, error) {
	return s.enf.AddGroupingPolicy(email, role)
}

func (s *RoleService) RevokeRole(email, role string) (bool, error) {
	return s.enf.RemoveGroupingPolicy(email, role)
}

func (s *RoleService) GrantPermission(role, obj, act string) (bool, error) {
	return s.enf.AddPolicy(role, obj, act)
}

func (s *RoleService) RevokePermission(role, obj, act string) (bool, error) {
	return s.enf.RemovePolicy(role, obj, act)
}

func (s *RoleService) GetRoleForUser(email string) ([]string, error) {
	return s.enf.GetRolesForUser(email)
}
