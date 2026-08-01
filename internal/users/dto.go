package users

type CreateUserDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
	Image    string `json:"image"`
	OfficeID uint64 `json:"office_id" binding:"required"`
	Role     string `json:"role" binding:"required,role"`
}

type UpdateUserDTO struct {
	Id       uint64 `json:"id"`
	Email    string `json:"email" binding:"omitempty,email"`
	Name     string `json:"name"`
	Password string `json:"password" binding:"omitempty,min=8"`
	Image    string `json:"image"`
	OfficeID uint64 `json:"office_id"`
	Role     string `json:"role" binding:"omitempty,role"`
}
