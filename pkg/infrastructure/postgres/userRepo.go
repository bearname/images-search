package postgres

import (
	"github.com/col3name/images-search/pkg/domain/user"
	"github.com/jackc/pgx"
)

type UserRepository struct {
	connPool *pgx.ConnPool
}

func NewUserRepository(connPool *pgx.ConnPool) *UserRepository {
	u := new(UserRepository)
	u.connPool = connPool
	return u
}

func (r *UserRepository) CreateUser(username string, password []byte, role user.Role) error {
	sqlQuery := r.getCreateUserSql()
	query, err := r.connPool.Query(sqlQuery, username, string(password), role)
	if err != nil {
		return err
	}

	defer query.Close()

	return nil
}

func (r *UserRepository) getCreateUserSql() string {
	return "INSERT INTO users (username, passwd, user_role) VALUES ($1, $2, $3);"
}

func (r *UserRepository) FindByUserName(username string) (user.User, error) {
	var dbUser user.User

	row := r.connPool.QueryRow(r.getFindUserSql(), username)

	err := row.Scan(
		&dbUser.Id,
		&dbUser.Username,
		&dbUser.Password,
		&dbUser.Role,
	)

	return dbUser, err
}

func (r *UserRepository) getFindUserSql() string {
	return "SELECT id, username, passwd, user_role FROM users WHERE username = $1;"
}
