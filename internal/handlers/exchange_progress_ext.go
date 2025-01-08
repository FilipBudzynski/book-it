package handlers

import (
	"fmt"
	"reflect"

	"github.com/labstack/echo/v4"
)

type exchangeFormBinding struct {
	DesiredBookID string `form:"desired-book-id"`
	UserBookIDs   []string
}

func (e *exchangeFormBinding) bind(c echo.Context) error {
	val := reflect.ValueOf(e).Elem()
	typ := reflect.TypeOf(*e)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		formTag := field.Tag.Get("form")
		if formTag != "" {
			if value := c.FormValue(formTag); value != "" {
				val.Field(i).SetString(value)
			}
		}
	}
	e.UserBookIDs = []string{}
	for i := 0; i <= 4; i++ {
		fieldName := fmt.Sprintf("user-book-%d", i)
		if value := c.FormValue(fieldName); value != "" {
			e.UserBookIDs = append(e.UserBookIDs, value)
		}
	}
	return nil
}
