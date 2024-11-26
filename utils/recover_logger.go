package utils

import (
	"fmt"
	"runtime"

	"github.com/labstack/echo/v4"
)

// CustomRecoverMiddleware handles panics and logs the stack trace in plain-text format
func CustomRecoverMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if r := recover(); r != nil {
				stack := make([]byte, 4*1024) // Adjust size if needed
				length := runtime.Stack(stack, false)
				stackTrace := string(stack[:length])

				// Log stack trace in plain-text format
				fmt.Printf("PANIC RECOVERED: %v\n%s\n", r, stackTrace)

				// Send internal server error response
				// c.JSON(http.StatusInternalServerError, map[string]string{
				// 	"message": "Internal Server Error",
				// })
			}
		}()
		return next(c)
	}
}
