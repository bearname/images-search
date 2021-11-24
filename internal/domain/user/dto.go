package user

type FindUserDto struct {
	Username string `json:"usernameOrId"`
	Role     int    `json:"role"`
}
