package request

type AuthCredentials struct {
	Login    string `json:"login" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8,max=24"`
}
