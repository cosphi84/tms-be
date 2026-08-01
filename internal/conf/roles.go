package conf

import "slices"

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

// AllRoles = SATU-SATUNYA tempat daftar role didefinisikan. Semua tempat lain
// (validasi DTO, seeder, dsb) WAJIB baca dari sini, BUKAN hardcode ulang.
// Nambah/hapus role -> cukup edit slice ini, efeknya otomatis nyebar ke mana-mana.
var AllRoles = []RoleType{
	RoleSuperadmin,
	RoleAdminHQ,
	RoleSvcHead,
	RoleManagement,
	RoleAuditor,
	RoleAuditorCS,
	RoleTechnician,
}

// Valid cek apakah RoleType ini terdaftar di AllRoles.
func (r RoleType) Valid() bool {
	return slices.Contains(AllRoles, r)
}

type Key string

const (
	AuthContextKey Key = "auth"
)
