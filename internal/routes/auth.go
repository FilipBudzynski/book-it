package routes

import (
	"github.com/FilipBudzynski/book_it/internal/handlers"
	"github.com/labstack/echo/v4"
)

func RegisterAuthRoutes(app *echo.Echo, authHandler *handlers.AuthHandler) {
	group := app.Group("/auth")
	group.GET("/callback", authHandler.GetAuthCallbackFunc)
	group.GET("/", authHandler.GetAuthFunc)
	group.GET("/logout", authHandler.Logout)
}
