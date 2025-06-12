package utils

import (
    "github.com/gin-gonic/gin"
)

type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message,omitempty"`
    Code    string `json:"code,omitempty"`
}

func RespondWithError(c *gin.Context, code int, message string) {
    c.JSON(code, ErrorResponse{
        Error: message,
    })
}

func RespondWithDetailedError(c *gin.Context, code int, err string, message string, errorCode string) {
    c.JSON(code, ErrorResponse{
        Error:   err,
        Message: message,
        Code:    errorCode,
    })
}