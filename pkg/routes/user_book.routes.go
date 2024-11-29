package routes

import (
	"github.com/FilipBudzynski/book_it/pkg/handlers"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/labstack/echo/v4"
)

func RegisterUserBookRoutes(app *echo.Echo, h *handlers.UserBookHandler) {
	group := app.Group("/user-books")
	// middleware for protected routes
	group.Use(utils.CheckLoggedInMiddleware)
	// UserBook endpoints
	group.POST("/:book_id", h.Create)
	group.DELETE("/:book_id", h.Delete)
	group.GET("", h.List)
}
