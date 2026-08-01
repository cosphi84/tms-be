package auth

type LoginRequestDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginResponseDTO struct {
	AccessToken string      `json:"access_token"`
	User        interface{} `json:"user"`

	// RefreshToken SENGAJA json:"-" — gak boleh pernah muncul di response
	// body. Handler baca field ini buat di-set sebagai HttpOnly cookie,
	// tapi field ini tidak pernah ter-serialize ke JSON yang dikirim ke FE.
	RefreshToken string `json:"-"`
}

type RefreshTokenRequestDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
