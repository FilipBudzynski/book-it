package routes

import (
	"github.com/FilipBudzynski/book_it/internal/handlers"
	"github.com/labstack/echo/v4"
)

func RegisterBookRoutes(app *echo.Echo, h *handlers.BookHandler) {
	group := app.Group("/books")
	group.GET("", h.ListBooks)
	group.POST("", h.ListBooks)
}
