package v1

type RegisterUserDto struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

func NewRegisterUserDto(username string, password string) *RegisterUserDto {
	return &RegisterUserDto{
		Username: username,
		Password: password,
	}
}
