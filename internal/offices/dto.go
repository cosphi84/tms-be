package offices

type OfficeRequestDTO struct {
	Code string `json:"code" binding:"required,max=10"`
	Name string `json:"name" binding:"required,max=100"`
	Type string `json:"type" binding:"required,oneof=cabang sdss ssr sass tc"`

	// WAJIB diisi — office biasa (non-HQ) selalu punya parent.
	// HQ/root cuma bisa dibuat lewat CreateHQ (dari seeder), bukan endpoint ini.
	ParentID *uint64 `json:"parent_id" binding:"required"`
}

// HQRequestDTO SENGAJA gak punya ParentID — root gak punya parent, dan
// cuma boleh dibuat sekali lewat seeder (lihat catatan di service.go).
type HQRequestDTO struct {
	Code string `json:"code" binding:"required,max=10"`
	Name string `json:"name" binding:"required,max=100"`
	Type string `json:"type" binding:"required,oneof=hq"`
}

// OfficeTreeOption = bentuk data buat komponen select/combobox nested di FE
// (shadcn). Value/Label itu konvensi standar shadcn combobox, Children buat
// render submenu/grup bertingkat sesuai hierarki nested set.
type OfficeTreeOption struct {
	Value    string              `json:"value"` // office ID, string biar aman dipakai langsung sebagai HTML value
	Label    string              `json:"label"`
	Code     string              `json:"code"`
	Type     string              `json:"type"`
	Children []*OfficeTreeOption `json:"children,omitempty"`
}
