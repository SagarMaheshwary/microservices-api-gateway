package types

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
