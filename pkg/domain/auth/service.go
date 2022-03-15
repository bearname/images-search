package auth

type Service interface {
	CreateUser() (Token, error)
	Login(loginUserDto *Credentials) (Token, error)
	ValidateToken(authorizationHeader string) (string, error)
	RefreshToken(refreshTokenDto *RefreshTokenDto) (Token, error)
}
