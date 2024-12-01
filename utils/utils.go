package utils

import (
	"github.com/FilipBudzynski/book_it/pkg/models"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func RenderView(c echo.Context, cmp templ.Component) error {
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

func BookInUserBooks(bookID string, userBooks []models.UserBook) bool {
	for _, userBook := range userBooks {
		if userBook.BookID == bookID {
			return true
		}
	}
	return false
}
