package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/FilipBudzynski/book_it/pkg/models"
	"github.com/FilipBudzynski/book_it/pkg/services"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const (
	maxAge = 86400
	isProd = false
)

func UseAuth() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	hashingKey := os.Getenv("HASHING_KEY")

	store := sessions.NewCookieStore([]byte(hashingKey))
	store.MaxAge(maxAge)
	store.Options.HttpOnly = true
	store.Options.Secure = true

	gothic.Store = store

	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, "http://localhost:3000/auth/google/callback"),
	)
}

type AuthHandler struct {
	userService services.UserService
}

func NewAuthHandler(us services.UserService) *AuthHandler {
	return &AuthHandler{
		userService: us,
	}
}

func (a *AuthHandler) setProvider(c echo.Context) (http.ResponseWriter, *http.Request) {
	ctx := context.WithValue(c.Request().Context(), gothic.ProviderParamKey, c.Param("provider"))
	return c.Response().Writer, c.Request().WithContext(ctx)
}

func (a *AuthHandler) GetAuthCallbackFunc(c echo.Context) error {
	response, request := a.setProvider(c)

	gothUser, err := gothic.CompleteUserAuth(response, request)
	if err != nil {
		fmt.Fprintln(c.Response().Writer, err.Error())
		return err
	}

	// try to get user from db
	user, err := a.userService.GetByEmail(gothUser.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	// create new user
	if user == nil {
		user = &models.User{
			Username: gothUser.Name,
			Email:    gothUser.Email,
			GoogleId: gothUser.UserID,
		}
		if err = a.userService.Create(user); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	log.Printf("successfully logged-in user: %s", user.Username)
	// set cookie session
	session, _ := gothic.Store.New(c.Request(), "session")
	session.Values["userID"] = user.GoogleId
	session.Save(c.Request(), c.Response().Writer)

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (a *AuthHandler) GetAuthFunc(c echo.Context) error {
	response, request := a.setProvider(c)

	if gothUser, err := gothic.CompleteUserAuth(c.Response().Writer, c.Request()); err == nil {
		fmt.Fprintln(c.Response().Writer, gothUser)
	} else {
		gothic.BeginAuthHandler(response, request)
	}

	return nil
}

func (a *AuthHandler) Logout(c echo.Context) error {
    // remove cookie from server
	gothic.Logout(c.Response().Writer, c.Request())

	// remove cookie from user side
	session, _ := gothic.Store.Get(c.Request(), "session")
	session.Options.MaxAge = -1
	session.Save(c.Request(), c.Response().Writer)

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}
