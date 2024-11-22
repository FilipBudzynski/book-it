package routes

import (
	"github.com/FilipBudzynski/book_it/pkg/handlers"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/labstack/echo/v4"
)

func RegisterUserRoutes(app *echo.Echo, h *handlers.UserHandler) {
	group := app.Group("/users")
	// middleware - protected routes
	// group.Use()
	group.Use(utils.CheckLoggedInMiddleware)
	group.GET("", h.ListUsers)
    group.GET("/add", h.AddBook)
}
