package helper_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DummyInput struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type DummyValidationError struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func TestGetRootDir(t *testing.T) {
	root := helper.GetRootDir()
	assert.NotEmpty(t, root)
}

func TestPrepareResponse(t *testing.T) {
	message := "success"
	data := gin.H{"id": 1}
	res := helper.PrepareResponse(message, data)

	assert.Equal(t, message, res["message"])
	assert.Equal(t, data, res["data"])
}

func TestValidationErrorByTag(t *testing.T) {
	tests := []struct {
		tag      string
		field    string
		expected string
	}{
		{"required", "email", "email is a required"},
		{"email", "email", "email must be an email"},
		{"unknown", "field", ""},
	}

	for _, tt := range tests {
		got := helper.ValidationErrorByTag(tt.tag, tt.field)
		assert.Equal(t, tt.expected, got)
	}
}

func TestPrepareResponseFromValidationError(t *testing.T) {
	validate := validator.New()
	input := &DummyInput{}
	validationErr := validate.Struct(input)
	validationObj := &DummyValidationError{}

	res := helper.PrepareResponseFromValidationError(validationErr, validationObj)
	data := res["data"].(gin.H)
	errorMap := data["errors"].(map[string][]string)
	assert.Contains(t, errorMap, "name")
	assert.Contains(t, errorMap, "email")
	assert.Equal(t, constant.MessageBadRequest, res["message"])
}

func TestPrepareResponseFromGRPCError(t *testing.T) {
	// BAD REQUEST with JSON error
	validationErrors := &DummyValidationError{"invalid name", "invalid email"}
	payload, err := json.Marshal(validationErrors)
	require.NoError(t, err)

	grpcErr := status.Error(codes.InvalidArgument, string(payload))
	code, res := helper.PrepareResponseFromGRPCError(grpcErr, &DummyValidationError{})

	assert.Equal(t, http.StatusBadRequest, code)
	assert.Equal(t, validationErrors, res["data"].(gin.H)["errors"])

	// INTERNAL ERROR with no JSON payload
	grpcErr = status.Error(codes.Internal, constant.MessageInternalServerError)
	code, res = helper.PrepareResponseFromGRPCError(grpcErr, &DummyValidationError{})

	assert.Equal(t, http.StatusInternalServerError, code)
	assert.Equal(t, constant.MessageInternalServerError, res["message"])
	assert.Empty(t, res["data"].(gin.H))
}

func TestGRPCToHttpCode(t *testing.T) {
	tests := []struct {
		code     codes.Code
		expected int
	}{
		{codes.OK, http.StatusOK},
		{codes.InvalidArgument, http.StatusBadRequest},
		{codes.Unauthenticated, http.StatusUnauthorized},
		{codes.PermissionDenied, http.StatusForbidden},
		{codes.NotFound, http.StatusNotFound},
		{codes.AlreadyExists, http.StatusConflict},
		{codes.Internal, http.StatusInternalServerError},
		{codes.Unavailable, http.StatusServiceUnavailable},
		{codes.Unknown, http.StatusInternalServerError}, // fallback
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, helper.GRPCToHttpCode(tt.code))
	}
}

func TestHTTPCodeToMessage(t *testing.T) {
	tests := []struct {
		code     int
		expected string
	}{
		{http.StatusOK, constant.MessageOK},
		{http.StatusBadRequest, constant.MessageBadRequest},
		{http.StatusUnauthorized, constant.MessageUnauthorized},
		{http.StatusForbidden, constant.MessageForbidden},
		{http.StatusNotFound, constant.MessageNotFound},
		{http.StatusConflict, constant.MessageConflict},
		{http.StatusInternalServerError, constant.MessageInternalServerError},
		{http.StatusServiceUnavailable, constant.MessageServiceUnavailable},
		{418, constant.MessageInternalServerError}, // unknown/fallback
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, helper.HTTPCodeToMessage(tt.code))
	}
}
