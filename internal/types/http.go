package types

type AuthorizationHeader struct {
	Token string `header:"authorization" binding:"required"`
}

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

type VerifyTokenValidationError struct {
	//
}

type UploadedWebhookInput struct {
	VideoId     string `json:"video_id" binding:"required"`
	ThumbnailId string `json:"thumbnail_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type UploadedWebhookValidationError struct {
	VideoId     []string `json:"video_id"`
	ThumbnailId []string `json:"thumbnail_id"`
	Title       []string `json:"title"`
	Description []string `json:"description"`
}
