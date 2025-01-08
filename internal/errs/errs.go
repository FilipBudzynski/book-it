package errs

import (
	"net/http"

	"github.com/FilipBudzynski/book_it/internal/toast"
	"github.com/labstack/echo/v4"
)

var (
	HttpErrorBadRequest = func(err error) *echo.HTTPError {
		toast := toast.Warning(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, toast).
			SetInternal(toast)
	}

	HttpErrorUnauthorized = func(err error) *echo.HTTPError {
		toast := toast.Warning(err.Error())
		return echo.NewHTTPError(http.StatusUnauthorized, toast).
			SetInternal(toast)
	}

	HttpErrorConflict = func(err error) *echo.HTTPError {
		toast := toast.Warning(err.Error())
		return echo.NewHTTPError(http.StatusConflict, toast).
			SetInternal(toast)
	}

	HttpErrorInternalServerError = func(err error) *echo.HTTPError {
		toast := toast.Danger(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, toast).
			SetInternal(toast)
	}
)
