package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/FilipBudzynski/book_it/pkg/models"
	"github.com/FilipBudzynski/book_it/pkg/services"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

type AuthHandler struct {
	userService services.UserService
}

func NewAuthHandler(us services.UserService) *AuthHandler {
	return &AuthHandler{
		userService: us,
	}
}

// setProvider is a helper function that sets Request context to contain value "provider", from url path ":provider"
// returns responseWriter and altered request
func setProvider(c echo.Context) (http.ResponseWriter, *http.Request) {
	ctx := context.WithValue(c.Request().Context(), gothic.ProviderParamKey, c.Param("provider"))
	return c.Response().Writer, c.Request().WithContext(ctx)
}

func (a *AuthHandler) GetAuthCallbackFunc(c echo.Context) error {
	responseWriter, request := setProvider(c)

	if request.URL.Query().Get("code") == "" {
		log.Println("user has canceled authentication")
		gothic.Logout(responseWriter, request)
		return c.Redirect(http.StatusFound, "/")
	}

	gothUser, err := gothic.CompleteUserAuth(responseWriter, request)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// try to get user from db
	user, err := a.userService.GetByEmail(gothUser.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	// create new user
	if user == nil {
		user = &models.User{
			Username: gothUser.NickName,
			Email:    gothUser.Email,
			GoogleId: gothUser.UserID,
		}
		if err = a.userService.Create(user); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}
	log.Printf("successfully logged-in user: %s", user.Username)

	// set cookie session
	err = utils.SetSessionUserID(responseWriter, request, gothUser.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.Redirect(http.StatusFound, "/")
}

func (a *AuthHandler) GetAuthFunc(c echo.Context) error {
	responseWriter, request := setProvider(c)
	gothic.BeginAuthHandler(responseWriter, request)
	return nil
}

func (a *AuthHandler) Logout(c echo.Context) error {
	r := c.Request()
	w := c.Response().Writer

	// remove cookie
	err := utils.RemoveCookieSession(w, r)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.Redirect(http.StatusFound, "/")
}
