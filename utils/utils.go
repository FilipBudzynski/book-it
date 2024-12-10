package utils

import (
	"github.com/FilipBudzynski/book_it/cmd/web"
	"github.com/FilipBudzynski/book_it/pkg/models"
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

func HTMXMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Check if the request is an HTMX request
		if c.Request().Header.Get("HX-Request") == "false" {
			// If not HTMX, render the base layout and pass the requested route as content
			RenderView(c, web.Base())
		}
		return next(c)
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
