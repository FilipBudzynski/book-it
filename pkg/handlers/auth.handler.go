package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const (
	maxAge = 86400 * 30
	isProd = false
)

func NewAuth() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	// hashingKey := os.Getenv("HASHING_KEY")
	//
	// store := sessions.NewCookieStore([]byte(hashingKey))
	// store.MaxAge(maxAge)
	// store.Options.Path = "/"
	// store.Options.HttpOnly = true
	// store.Options.Secure = isProd
	//
	// gothic.Store = store

	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, "http://localhost:3000/auth/google/callback"),
	)
}

func GetAuthCallbackFunc(c echo.Context) error {
	ctx := context.WithValue(c.Request().Context(), gothic.ProviderParamKey, c.Param("provider"))
	response := c.Response().Writer
	request := c.Request().WithContext(ctx)

	user, err := gothic.CompleteUserAuth(response, request)
	if err != nil {
		fmt.Fprintln(c.Response().Writer, err.Error())
		return err
	} else {
		fmt.Fprintln(c.Response().Writer, user)
	}

	return c.Redirect(http.StatusFound, "/")
}

func GetAuthFunc(c echo.Context) error {
	ctx := context.WithValue(c.Request().Context(), gothic.ProviderParamKey, c.Param("provider"))
	response := c.Response().Writer
	request := c.Request().WithContext(ctx)

	// if gothUser, err := gothic.CompleteUserAuth(c.Response().Writer, c.Request()); err == nil {
	// 	fmt.Fprintln(c.Response().Writer, gothUser)
	// } else {
	gothic.BeginAuthHandler(response, request)
	// }

	return nil
}

func Logout(c echo.Context) error {
	gothic.Logout(c.Response().Writer, c.Request())
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}
