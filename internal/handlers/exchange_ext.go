package handlers

import (
	"fmt"
	"reflect"

	"github.com/labstack/echo/v4"
)

var ExchangeDeclineAlertMessage = func(title, user string) string {
	return fmt.Sprintf("Your exchange request for '%s' was declined by %s.", title, user)
}

var ExchangeAcceptedAlertMessage = func(title, user string) string {
	return fmt.Sprintf("Your exchange request for '%s' was ACCEPTED by %s.", title, user)
}

type exchangeFormBinding struct {
	DesiredBookID string `form:"desired-book-id"`
	UserBookIDs   []string
	Latitude      string `form:"latitude"`
	Longitude     string `form:"longitude"`
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
		fieldName := fmt.Sprintf("offered-book-%d", i)
		if value := c.FormValue(fieldName); value != "" {
			e.UserBookIDs = append(e.UserBookIDs, value)
		}
	}
	return nil
}
