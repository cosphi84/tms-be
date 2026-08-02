package sloc

type CreateSlocDTO struct {
	Code     string `json:"code" binding:"required,max=20"`
	Name     string `json:"name" binding:"required,max=100"`
	OfficeID uint64 `json:"office_id" binding:"required"`
	Member   string `json:"member" binding:"required,slocMember"`
}

type UpdateSlocDTO struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	OfficeID uint64 `json:"office_id"`
	Member   string `json:"member" binding:"omitempty,slocMember"`
}

type SlocOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}
