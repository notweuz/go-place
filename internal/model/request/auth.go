package request

type AuthCredentials struct {
	Login    string `json:"login" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}
