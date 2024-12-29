package utils

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/FilipBudzynski/book_it/cmd/web"
	"github.com/FilipBudzynski/book_it/internal/models"
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

func BookInUserBooks(bookID string, userBooks []*models.UserBook) bool {
	for _, userBook := range userBooks {
		if userBook.BookID == bookID {
			return true
		}
	}
	return false
}

// StructNameToSnakeCase takes a struct and returns its name in snake_case
func StructNameToSnakeCase[T any]() string {
	t := reflect.TypeOf((*T)(nil))

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return ""
	}

	structName := t.Name()

	var snakeCaseName strings.Builder
	for i, r := range structName {
		if unicode.IsUpper(r) {
			if i > 0 {
				snakeCaseName.WriteRune('_')
			}
			snakeCaseName.WriteRune(unicode.ToLower(r))
		} else {
			snakeCaseName.WriteRune(r)
		}
	}

	return snakeCaseName.String()
}

func StructName[T any]() string {
	t := reflect.TypeOf((*T)(nil))

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return ""
	}

	return t.Name()
}
