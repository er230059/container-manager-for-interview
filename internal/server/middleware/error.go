package middleware

import (
	"errors"
	"net/http"

	customErr "container-manager/internal/errors"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var customErr *customErr.CustomError
			if errors.As(err, &customErr) {
				c.JSON(customErr.Status, gin.H{"error": customErr.Message})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
		}
	}
}
