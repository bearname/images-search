package user

type Role int

const (
    Admin Role = iota
    General
)

func (q Role) String() string {
    return [...]string{"Admin", "General"}[q]
}

func (q Role) Values() int {
    return [...]int{0, 1}[q]
}

type User struct {
    Id       int    `json:"id"`
    Username string `json:"username"`
    Password []byte `json:"password"`
    Role     Role   `json:"role"`
}
