package constant

// response messages
const (
	MessageOK                  = "Success"
	MessageCreated             = "Created New Resource"
	MessageBadRequest          = "Bad Request"
	MessageUnauthorized        = "Unauthorized"
	MessageForbidden           = "Forbidden"
	MessageNotFound            = "Resource Not Found"
	MessageConflict            = "Conflict"
	MessageInternalServerError = "Internal Server Error"
	MessageServiceUnavailable  = "Service Unavailable"
)

const (
	MessageServicesUnhealthy = "Some services are not available!"
	MessageServicesHealthy   = "All services are healthy!"
)

const (
	HealthStatusHealthy  = "healthy"
	HealthStatusDegraded = "degraded"
)

// gRPC metadata headers
const (
	GRPCHeaderAuthorization = "authorization"
	GRPCHeaderUserId        = "x-user-id"
)

const (
	AuthUser = "user"
)

const ServiceName = "API Gateway"

const ExitFailure = 1
