package utils

import (
	"net/http"
	"strconv"
	"time"

	"github.com/FilipBudzynski/book_it/cmd/web"
	webError "github.com/FilipBudzynski/book_it/cmd/web/error_pages"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func RenderView(c echo.Context, cmp templ.Component) error {
	requestContext := c.Request().Context()
	responseWriter := c.Response().Writer

	if c.Request().Header.Get("HX-Request") == "true" {
		return cmp.Render(requestContext, responseWriter)
	} else {
		ctx := templ.WithChildren(requestContext, cmp)
		return web.Base().Render(ctx, responseWriter)
	}
}

func ErrorPagesMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			if httpErr, ok := err.(*echo.HTTPError); ok {
				c.Response().WriteHeader(httpErr.Code)
				c.Response().Header().Set("HX-Retarget", "#content-container")
				switch httpErr.Code {
				case http.StatusNotFound:
					return RenderView(c, webError.ErrorPage404())
				case http.StatusUnauthorized:
					return RenderView(c, webError.ErrorPage401())
				}
			}
		}
		return err
	}
}

func TodaysDate() time.Time {
	now := time.Now()
	year, month, day := now.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return today
}

func ParseStringToUint(s string) (uint, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return uint(i), nil
}
