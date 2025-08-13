package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetRootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))

	return filepath.Dir(d)
}

func PrepareResponse(message string, data any) gin.H {
	return gin.H{
		"message": message,
		"data":    data,
	}
}

func PrepareResponseFromValidationError(err error, obj any) gin.H {
	errorsMap := map[string][]string{}

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			f, _ := reflect.TypeOf(obj).Elem().FieldByName(e.Field())
			field, _ := f.Tag.Lookup("json")
			errorsMap[field] = []string{ValidationErrorByTag(e.Tag(), field)}
		}

		// Add empty slices for fields without errors to keep structure consistent
		fields := reflect.VisibleFields(reflect.Indirect(reflect.ValueOf(obj)).Type())
		for _, field := range fields {
			t, _ := field.Tag.Lookup("json")
			if _, ok := errorsMap[t]; !ok {
				errorsMap[t] = []string{}
			}
		}

		return PrepareResponse(constant.MessageBadRequest, gin.H{
			"errors": errorsMap,
		})
	}

	return PrepareResponse(constant.MessageBadRequest, gin.H{
		"errors": errorsMap,
	})
}

func ValidationErrorByTag(tag string, field string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "invalid":
		return fmt.Sprintf("%s is invalid", field)
	case "email":
		return fmt.Sprintf("%s must be an email", field)
	}
	return ""
}

func PrepareResponseFromGRPCError(err error, obj any) (int, gin.H) {
	e, _ := status.FromError(err)

	s := GRPCToHttpCode(e.Code())
	data := gin.H{}

	if s == http.StatusBadRequest {
		json.Unmarshal([]byte(e.Message()), &obj)
		data["errors"] = obj
	}

	res := PrepareResponse(HTTPCodeToMessage(s), data)

	return s, res
}

func GRPCToHttpCode(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

func HTTPCodeToMessage(code int) string {
	switch code {
	case http.StatusOK:
		return constant.MessageOK
	case http.StatusBadRequest:
		return constant.MessageBadRequest
	case http.StatusUnauthorized:
		return constant.MessageUnauthorized
	case http.StatusForbidden:
		return constant.MessageForbidden
	case http.StatusNotFound:
		return constant.MessageNotFound
	case http.StatusConflict:
		return constant.MessageConflict
	case http.StatusInternalServerError:
		return constant.MessageInternalServerError
	case http.StatusServiceUnavailable:
		return constant.MessageServiceUnavailable
	default:
		return constant.MessageInternalServerError
	}
}

func GetEnv(key string, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultVal
}

func GetEnvInt(key string, defaultVal int) int {
	if val, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return val
	}

	return defaultVal
}

func GetEnvDurationSeconds(key string, defaultVal time.Duration) time.Duration {
	if val, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return time.Duration(val) * time.Second
	}

	return defaultVal * time.Second
}
