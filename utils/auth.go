package utils

import (
	"encoding/gob"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

const (
	SessionName = "_user_session"
	userIDKey   = "user_id"
	maxAge      = 1800
)

type UserSession struct {
	UserID       string
	AccessToken  string
	RefreshToken string
}

func init() {
	gob.Register(UserSession{})
}

// setSessionValue stores a value in the current user session.
func setSessionValue(w http.ResponseWriter, r *http.Request, key, value any) error {
	session, _ := gothic.Store.New(r, SessionName)
	session.Options.MaxAge = maxAge
	session.Values[key] = value
	return session.Save(r, w)
}

func IsUserLoggedIn(r *http.Request) bool {
	session, _ := gothic.Store.Get(r, SessionName)
	_, ok := session.Values[userIDKey]
	return ok
}

func GetUserSessionFromStore(r *http.Request) (UserSession, error) {
	session, err := gothic.Store.Get(r, SessionName)
	if value, ok := session.Values["userSession"]; ok {
		return value.(UserSession), nil
	}
	return UserSession{}, fmt.Errorf("eser session not found %s", err)
}

// GetUserIDFromSession is an abstraction over GetFromSession to retrive userID from session
func GetUserIDFromSession(r *http.Request) (string, error) {
	user, err := GetUserSessionFromStore(r)
	if err != nil {
		return "", err
	}
	return user.UserID, nil
}

func SetUserSession(w http.ResponseWriter, r *http.Request, userSession UserSession) error {
	return setSessionValue(w, r, "userSession", userSession)
}

// SetSessionValue stores a value in the current user session.
func SetSessionValueMap(w http.ResponseWriter, r *http.Request, values map[any]any) error {
	session, _ := gothic.Store.New(r, SessionName)
	session.Values = values
	return session.Save(r, w)
}

// GetFromSession is an abstraction over SetSessionValue to store userID in session.
func SetSessionUserID(w http.ResponseWriter, r *http.Request, id string) error {
	return setSessionValue(w, r, userIDKey, id)
}

func RemoveCookieSession(w http.ResponseWriter, r *http.Request) error {
	session, _ := gothic.Store.Get(r, SessionName)
	session.Options.MaxAge = -1
	// session.Values = make(map[any]any)
	return session.Save(r, w)
}

func CheckLoggedInMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := GetUserSessionFromStore(c.Request())
		if err != nil {
			return echo.NewHTTPError(echo.ErrUnauthorized.Code, "Unauthorized")
		}
		return next(c)
	}
}

func RefreshSessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, _ := gothic.Store.Get(c.Request(), SessionName)

		if !session.IsNew {
			session.Options.MaxAge = 1800
			if err := session.Save(c.Request(), c.Response().Writer); err != nil {
				return err
			}
		}

		return next(c)
	}
}
