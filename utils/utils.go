package utils

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

const (
	sessionName        = "_user_session"
	userIDInSessionKey = "_user_session"
)

func RenderView(c echo.Context, cmp templ.Component) error {
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

func IsUserLoggedIn(r *http.Request) bool {
	session, _ := gothic.Store.Get(r, sessionName)
	_, ok := session.Values[userIDInSessionKey]
	return ok
}

// GetFromSession retrives a previously-stored value from the session.
func GetFromSession(r *http.Request, key string) (string, error) {
	session, err := gothic.Store.Get(r, sessionName)
	if value, ok := session.Values[key]; ok {
		return value.(string), nil
	}
	return "", err
}

// GetSessionUserID is an abstraction over GetFromSession to retrive userID from session
func GetSessionUserID(r *http.Request) (string, error) {
	return GetFromSession(r, userIDInSessionKey)
}

// SetSessionValue stores a value in the current user session.
func SetSessionValue(w http.ResponseWriter, r *http.Request, key, value any) error {
	session, _ := gothic.Store.New(r, sessionName)
	session.Values[key] = value
	return session.Save(r, w)
}

// GetFromSession is an abstraction over SetSessionValue to store userID in session.
func SetSessionUserID(w http.ResponseWriter, r *http.Request, id string) error {
	return SetSessionValue(w, r, userIDInSessionKey, id)
}

func RemoveCookieSession(w http.ResponseWriter, r *http.Request) error {
	session, _ := gothic.Store.Get(r, sessionName)
	session.Options.MaxAge = -1
	session.Values = make(map[interface{}]interface{})
	return session.Save(r, w)
}
