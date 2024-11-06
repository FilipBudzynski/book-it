package handlers

import (
	"github.com/FilipBudzynski/book_it/cmd/web"
	"github.com/FilipBudzynski/book_it/pkg/entities"
	"github.com/FilipBudzynski/book_it/pkg/services"
	"github.com/labstack/echo/v4"
)

func NewUserHandler(us services.User) *UserHandler {
	return &UserHandler{
		userService: us,
	}
}

type UserHandler struct {
	userService services.User
}

func (u *UserHandler) CreateUser(c echo.Context) error {
	user := new(entities.User)

	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, err)
	}

	if err := u.userService.Create(user); err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, err)
	}

	return web.AppendUsersList(*user).Render(c.Request().Context(), c.Response().Writer)
}

func (u *UserHandler) ListUsers(c echo.Context) error {
	users, err := u.userService.GetAll()
	if err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, err)
	}

	return web.UserForm(users).Render(c.Request().Context(), c.Response().Writer)
}
