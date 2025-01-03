package utils

import (
	"fmt"

	"github.com/FilipBudzynski/book_it/internal/toast"
	"github.com/labstack/echo/v4"
)

func CustomErrorHandler(err error, c echo.Context) {
	te, ok := err.(toast.Toast)
	if !ok {
		fmt.Println(err)
		te = toast.Danger("there has been an unexpected error")
	}

	if te.Level != toast.SUCCESS {
		c.Response().Header().Set("HX-Reswap", "none")
	}

	te.SetHXTriggerHeader(c)
}
