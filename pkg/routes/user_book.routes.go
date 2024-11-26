package routes

import (
	"github.com/FilipBudzynski/book_it/pkg/handlers"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/labstack/echo/v4"
)

func RegisterUserBookRoutes(app *echo.Echo, h *handlers.UserBookHandler) {
	group := app.Group("/users")
	// middleware for protected routes
	group.Use(utils.CheckLoggedInMiddleware)
	// UserBook endpoints
	group.GET("/add", h.Create)
}
