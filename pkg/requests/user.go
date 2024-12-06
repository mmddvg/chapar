package requests

type User struct {
	Name     string `json:"name" validate:"required"`
	UserName string `json:"user_name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Login struct {
	UserName string `json:"user_name" validate:"required"`
	Password string `json:"password" validate:"required"`
}
