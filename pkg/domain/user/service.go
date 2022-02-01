package user

type Service interface {
	Find(usernameOrId string) (*FindUserDto, error)
}
