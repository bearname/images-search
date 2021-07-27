package postgres

import (
    "aws_rekognition_demo/internal/domain/user"
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
    sqlQuery := "INSERT INTO users (username, passwd, user_role) VALUES ($1, $2, $3);"
    query, err := r.connPool.Query(sqlQuery, username, string(password), role)
    if err != nil {
        return err
    }

    defer query.Close()

    return nil
}

func (r *UserRepository) FindByUserName(username string) (user.User, error) {
    var dbUser user.User

    row := r.connPool.QueryRow("SELECT id, username, passwd, user_role FROM users WHERE username = $1;", username)

    err := row.Scan(
        &dbUser.Id,
        &dbUser.Username,
        &dbUser.Password,
        &dbUser.Role,
    )

    return dbUser, err
}
