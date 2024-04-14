package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	cons "github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RootDir() string {
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
	errors := map[string][]string{}

	for _, e := range err.(validator.ValidationErrors) {
		f, _ := reflect.TypeOf(obj).Elem().FieldByName(e.Field())
		field, _ := f.Tag.Lookup("json")

		errors[field] = []string{ValidationErrorByTag(e.Tag(), field)}
	}

	fields := reflect.VisibleFields(reflect.Indirect(reflect.ValueOf(obj)).Type())

	//Set non-error key/value pair as empty slice to
	//keep "errors" field consistent with grpc response.
	for _, field := range fields {
		t, _ := field.Tag.Lookup("json")

		if _, ok := errors[t]; !ok {
			errors[t] = []string{}
		}
	}

	return PrepareResponse(cons.MessageBadRequest, gin.H{
		"errors": errors,
	})
}

func ValidationErrorByTag(tag string, field string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s is a required", field)
	case "email":
		return fmt.Sprintf("%s must be an email", field)
	}
	return ""
}

func PrepareResponseFromgrpcError(err error, obj any) (int, gin.H) {
	e, _ := status.FromError(err)

	status := GRPCTohttpCode(e.Code())
	data := gin.H{}

	if status == http.StatusBadRequest {
		json.Unmarshal([]byte(e.Message()), &obj)
		data["errors"] = obj
	}

	response := PrepareResponse(HTTPCodeToMessage(status), data)

	return status, response
}

func GRPCTohttpCode(code codes.Code) int {
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
	case codes.Internal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func HTTPCodeToMessage(code int) string {
	switch code {
	case http.StatusOK:
		return cons.MessageOK
	case http.StatusBadRequest:
		return cons.MessageBadRequest
	case http.StatusUnauthorized:
		return cons.MessageUnauthorized
	case http.StatusForbidden:
		return cons.MessageForbidden
	case http.StatusNotFound:
		return cons.MessageNotFound
	case http.StatusInternalServerError:
		return cons.MessageInternalServerError
	default:
		return cons.MessageInternalServerError
	}
}
