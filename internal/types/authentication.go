package types

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterValidationError struct {
	Name     []string `json:"name"`
	Email    []string `json:"email"`
	Password []string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginValidationError struct {
	Email    []string `json:"email"`
	Password []string `json:"password"`
}

type LogoutValidationError struct {
	//
}

type AuthorizationHeader struct {
	Token string `header:"authorization" binding:"required"`
}
