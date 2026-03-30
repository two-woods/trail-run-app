package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func RespondSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
	})
}

func RespondCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    data,
	})
}

func RespondAccepted(c *gin.Context, data interface{}) {
	c.JSON(http.StatusAccepted, APIResponse{
		Success: true,
		Data:    data,
	})
}

func RespondError(c *gin.Context, status int, code, message string) {
	c.JSON(status, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	})
}

func RespondValidationError(c *gin.Context, message string) {
	RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", message)
}

func RespondUnauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized"
	}
	RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", message)
}

func RespondForbidden(c *gin.Context, message string) {
	if message == "" {
		message = "Forbidden"
	}
	RespondError(c, http.StatusForbidden, "FORBIDDEN", message)
}

func RespondNotFound(c *gin.Context, message string) {
	if message == "" {
		message = "Resource not found"
	}
	RespondError(c, http.StatusNotFound, "NOT_FOUND", message)
}

func RespondInternalError(c *gin.Context, message string) {
	if message == "" {
		message = "Internal server error"
	}
	RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", message)
}
