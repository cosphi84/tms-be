package auth

import "tms-be/internal/users"

type AuthService interface {
	Login()
	Refresh()
}

type authService struct {
	userRepo *users.Repository
	roleSvc  *RoleService
}
