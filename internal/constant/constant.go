package constant

// response messages
const (
	MessageOK                  = "Success"
	MessageCreated             = "Created New Resource"
	MessageBadRequest          = "Bad Request"
	MessageUnauthorized        = "Unauthorized"
	MessageForbidden           = "Forbidden"
	MessageNotFound            = "Resource Not Found"
	MessageInternalServerError = "Internal Server Error"
)

// gRPC metadata headers
const (
	HeaderAuthorization = "authorization"
	HeaderUserId        = "x-user-id"
)

const (
	AuthUser = "user"
)
