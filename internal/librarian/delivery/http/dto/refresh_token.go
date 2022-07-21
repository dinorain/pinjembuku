package dto

type LibrarianRefreshTokenDto struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type LibrarianRefreshTokenResponseDto struct {
	AccessToken  string `json:"access_token" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}
