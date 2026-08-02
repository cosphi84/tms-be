package officeleaders

const DateLayout = "2006-01-02"

type AssignLeaderDTO struct {
	OfficeID  uint64 `json:"office_id" binding:"required"`
	UserID    uint64 `json:"user_id" binding:"required"`
	StartDate string `json:"start_date" binding:"required"` // format: YYYY-MM-DD
}

type EndTermDTO struct {
	EndDate string `json:"end_date" binding:"required"` // format: YYYY-MM-DD
}

// UpdateOfficeLeaderDTO buat koreksi tanggal doang (typo input, dsb) --
// BUKAN buat ganti office_id/user_id. Ganti office/user itu artinya
// "assignment baru", bukan "update assignment lama" -- pakai endpoint
// Assign lagi kalau itu maksudnya.
type UpdateOfficeLeaderDTO struct {
	StartDate string `json:"start_date"` // optional, kosong = gak diubah
	EndDate   string `json:"end_date"`   // optional, kosong = gak diubah
}
