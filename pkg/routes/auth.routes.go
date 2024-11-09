package routes

import (
	"github.com/FilipBudzynski/book_it/pkg/handlers"
	"github.com/labstack/echo/v4"
)

func RegisterAuthRoutes(app *echo.Echo, authHandler *handlers.AuthHandler) {
	group := app.Group("/auth")
	group.GET("/:provider/callback", authHandler.GetAuthCallbackFunc)
	group.GET("/:provider", authHandler.GetAuthFunc)
	group.GET("/:provider/logout", authHandler.Logout)
}
