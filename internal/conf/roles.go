package conf

type RoleType string

const (
	RoleSuperadmin RoleType = "superadmin"
	RoleAdminHQ    RoleType = "admin_hq"
	RoleSvcHead    RoleType = "service_head"
	RoleManagement RoleType = "management"
	RoleAuditor    RoleType = "auditor"
	RoleAuditorCS  RoleType = "auditor_cs"
	RoleTechnician RoleType = "technician"
)

type Key string

const (
	AuthContextKey Key = "auth"
)
