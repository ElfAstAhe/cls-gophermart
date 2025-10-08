package model

type User struct {
	ID       string `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Disabled bool   `db:"disabled"`

	Account *Account
}

func NewUser(username string, password string) *User {
	return &User{
		Username: username,
		Password: password,
		Disabled: false,
	}
}
