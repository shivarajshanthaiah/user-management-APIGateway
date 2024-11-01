package dto

type User struct {
	UserName string `json:"user_name" validate:"required,user_name"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,phone"`
}

type MethodRequest struct {
	Method   int `json:"method"`
	WaitTime int `json:"waitTime"`
}
