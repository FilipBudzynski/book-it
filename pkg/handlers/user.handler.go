package handlers

import (
	"fmt"

	"github.com/FilipBudzynski/book_it/cmd/web"
	"github.com/FilipBudzynski/book_it/pkg/models"
	"github.com/FilipBudzynski/book_it/pkg/services"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(us services.UserService) *UserHandler {
	return &UserHandler{
		userService: us,
	}
}

func (u *UserHandler) RegisterUser(c echo.Context) error {
	user := new(models.User)

	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, err)
	}

	if err := u.userService.Register(user); err != nil {
		return fmt.Errorf("registration failed, err %v", err) // echo.NewHTTPError(echo.ErrInternalServerError.Code, err)
	}

	fmt.Fprintln(c.Response().Writer, "User registered successfully")
	return nil

	// return web.AppendUsersList(*user).Render(c.Request().Context(), c.Response().Writer)
}

func (u *UserHandler) ListUsers(c echo.Context) error {
	users, err := u.userService.GetAll()
	if err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, err)
	}

	return web.UserForm(users).Render(c.Request().Context(), c.Response().Writer)
}

func (u *UserHandler) View(c echo.Context, cmp templ.Component) error {
	// c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}
