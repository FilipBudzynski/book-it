package handlers

import (
	"context"
	"net/http"

	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
	"gorm.io/gorm"
)

type AuthHandler struct {
	userService UserService
}

func NewAuthHandler(us UserService) *AuthHandler {
	return &AuthHandler{
		userService: us,
	}
}

func (h *AuthHandler) RegisterRoutes(app *echo.Echo) {
	group := app.Group("/auth")
	group.GET("/callback", h.GetAuthCallbackFunc)
	group.GET("", h.GetAuthFunc)
	group.GET("/logout", h.Logout)
}

// setProvider is a helper function that sets Request context to contain value "provider", from url path ":provider"
// returns responseWriter and altered request
func setProvider(c echo.Context) (http.ResponseWriter, *http.Request) {
	ctx := context.WithValue(c.Request().Context(), gothic.ProviderParamKey, c.QueryParam("provider"))
	return c.Response().Writer, c.Request().WithContext(ctx)
}

func (a *AuthHandler) GetAuthCallbackFunc(c echo.Context) error {
	responseWriter, request := setProvider(c)

	if request.URL.Query().Get("code") == "" {
		if err := gothic.Logout(responseWriter, request); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		c.Logger().Printf("user has canceled authentication")
		return c.Redirect(http.StatusFound, "/")
	}

	gothUser, err := gothic.CompleteUserAuth(responseWriter, request)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// try to get user from db
	// TODO: get by id
	var user *models.User

	user, err = a.userService.GetByEmail(gothUser.Email)
	// TODO: error should be wrapped
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			user = &models.User{
				Username: gothUser.Name,
				Email:    gothUser.Email,
				GoogleId: gothUser.UserID,
                AvatarURL: gothUser.AvatarURL,
			}
			if err = a.userService.Create(user); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	// set user cookie session
	err = utils.SetUserSession(responseWriter, request, utils.UserSession{
		UserID:       gothUser.UserID,
		UserEmail:    gothUser.Email,
		AccessToken:  gothUser.AccessToken,
		RefreshToken: gothUser.RefreshToken,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	c.Logger().Printf("successfully logged-in user: %s", user.Username)

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
