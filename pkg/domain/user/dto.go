package user

type FindUserDto struct {
	Id       int
	Username string `json:"usernameOrId"`
	Role     int    `json:"role"`
}
