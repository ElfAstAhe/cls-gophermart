package v1

type LoginDto struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

func NewLoginDto(username string, password string) *LoginDto {
	return &LoginDto{
		Username: username,
		Password: password,
	}
}
