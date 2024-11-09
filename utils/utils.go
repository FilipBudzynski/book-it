package utils

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

func RenderView(c echo.Context, cmp templ.Component) error {
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

func IsUserLoggedIn(r *http.Request, sessionName string) bool {
	session, _ := gothic.Store.Get(r, sessionName)
	_, ok := session.Values["userID"]
	return ok
}

func GetUserIDFromSession(r *http.Request, sessionName string) string {
	session, _ := gothic.Store.Get(r, sessionName)
	if userId, ok := session.Values["userID"]; ok {
		return userId.(string)
	}
	return ""
}
