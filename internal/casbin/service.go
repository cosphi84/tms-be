package casbin

import "github.com/casbin/casbin/v2"

type Service struct {
	enforcer *casbin.Enforcer
}

func New(enforcer *casbin.Enforcer) *Service {
	if enforcer == nil {
		panic("casbin enforcer is nil")
	}
	return &Service{enforcer: enforcer}
}

// Enforce cek apakah `sub` (user_id) boleh melakukan `act` terhadap `obj`.
// Role resolution didelegasikan penuh ke Casbin lewat g policy.
func (s *Service) Enforce(sub, obj, act string) (bool, error) {
	return s.enforcer.Enforce(sub, obj, act)
}

// GrantRole assign role ke user (nulis ke g policy).
// userID di sini HARUS identifier yang sama dengan yang dipakai di JWT claims.
func (s *Service) GrantRole(userID, role string) (bool, error) {
	return s.enforcer.AddGroupingPolicy(userID, role)
}

func (s *Service) RevokeRole(userID, role string) (bool, error) {
	return s.enforcer.RemoveGroupingPolicy(userID, role)
}

// GrantPermission assign permission ke role (nulis ke p policy).
func (s *Service) GrantPermission(role, obj, act string) (bool, error) {
	return s.enforcer.AddPolicy(role, obj, act)
}

func (s *Service) RevokePermission(role, obj, act string) (bool, error) {
	return s.enforcer.RemovePolicy(role, obj, act)
}

// GetRolesForUser GetRolesForUser: daftar role milik user — dipakai buat ditampilkan
// (misal endpoint GET /me), BUKAN buat authorization check.
func (s *Service) GetRolesForUser(userID string) ([]string, error) {
	return s.enforcer.GetRolesForUser(userID)
}
