package category

type CreateCategoryDTO struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"required,max=255"`
}

type UpdateCategoryDTO struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"required,max=255"`
}

type CategoryOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}
