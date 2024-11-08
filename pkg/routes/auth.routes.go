package routes

import (
	"github.com/FilipBudzynski/book_it/pkg/handlers"
	"github.com/labstack/echo/v4"
)

func RegisterAuthRoutes(app *echo.Echo) {
	group := app.Group("/auth")
	group.GET("/:provider/callback", handlers.GetAuthCallbackFunc)
	group.GET("/:provider", handlers.GetAuthFunc)
	group.GET("/logout/:provider", handlers.Logout)
}
