package errs

import (
	"net/http"

	"github.com/FilipBudzynski/book_it/internal/toast"
	"github.com/labstack/echo/v4"
)

var (
	HttpErrorBadRequest = func(err error) *echo.HTTPError {
		return echo.NewHTTPError(http.StatusBadRequest, toast.Warning(err.Error()))
	}

	HttpErrorUnauthorized = func(err error) *echo.HTTPError {
		return echo.NewHTTPError(http.StatusUnauthorized, toast.Warning(err.Error()))
	}

	HttpErrorConflict = func(err error) *echo.HTTPError {
		return echo.NewHTTPError(http.StatusConflict, toast.Danger(err.Error()))
	}

	HttpErrorInternalServerError = func(err error) *echo.HTTPError {
		return echo.NewHTTPError(http.StatusInternalServerError, toast.Danger(err.Error()))
	}
)
