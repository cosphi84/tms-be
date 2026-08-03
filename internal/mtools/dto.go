package mtools

type CreateToolDTO struct {
	Code        string  `json:"code" binding:"required,max=20"`
	Name        string  `json:"name" binding:"required,max=150"`
	Brand       string  `json:"merk" binding:"required,max=100"`
	Type        string  `json:"type" binding:"required,max=100"`
	SerialNum   string  `json:"serial_num"`
	CategoryID  uint32  `json:"category_id" binding:"required"`
	GroupID     uint32  `json:"group_id" binding:"required,max=100"`
	PhotoID     uint64  `json:"photo_id" binding:"omitempty,uuid"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	UsagePeriod string  `json:"usage_period" binding:"required,usagePeriod"`
}

// UpdateToolDTO -- Billable & IsActive pakai pointer (*bool), biar bisa
// bedain "field ini gak dikirim sama sekali" vs "dikirim dengan value false".
// Kalau pakai bool biasa, gak ada cara bedain "gak mau ubah" dari "mau
// diubah jadi false" -- field bool zero-value SELALU false.
type UpdateToolDTO struct {
	Code        string  `json:"code" binding:"required,max=20"`
	Name        string  `json:"name" binding:"required,max=150"`
	Brand       string  `json:"merk" binding:"required,max=100"`
	Type        string  `json:"type" binding:"required,max=100"`
	SerialNum   string  `json:"serial_num"`
	CategoryID  uint32  `json:"category_id" binding:"required"`
	GroupID     uint32  `json:"group_id" binding:"required,max=100"`
	PhotoID     uint64  `json:"photo_id" binding:"omitempty,uuid"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	UsagePeriod string  `json:"usage_period" binding:"required,usagePeriod"`
	IsActive    *bool   `json:"is_active"`
}
