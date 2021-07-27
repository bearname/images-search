package auth

type RefreshTokenDto struct {
    Username string `json:"username"`
    Token    string `json:"refreshToken"`
}
