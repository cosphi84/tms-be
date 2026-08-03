package groups

type CreateGroupDTO struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"required,max=255"`
}

type UpdateGroupDTO struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"required,max=255"`
}

type GroupOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}
