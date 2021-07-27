package user

type Repository interface {
    FindByUserName(username string) (User, error)
    CreateUser(username string, password []byte, role Role) error
}