package utils

import (
	"github.com/FilipBudzynski/book_it/cmd/web"
	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func RenderView(c echo.Context, cmp templ.Component) error {
	if c.Request().Header.Get("HX-Request") == "true" {
		return cmp.Render(c.Request().Context(), c.Response().Writer)
	} else {
		ctx := templ.WithChildren(c.Request().Context(), cmp)
		return web.Base().Render(ctx, c.Response().Writer)
	}
}

func BookInUserBooks(bookID string, userBooks []*models.UserBook) bool {
	for _, userBook := range userBooks {
		if userBook.BookID == bookID {
			return true
		}
	}
	return false
}
