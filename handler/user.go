package handler

import (
	"book_it/cmd/web"
	"book_it/service"
	"fmt"

	"github.com/labstack/echo/v4"
)

func NewUserHandler(us *service.UserService) *UserHandler {
	return &UserHandler{
		UserService: us,
	}
}

type UserHandler struct {
	UserService *service.UserService
}

func (u *UserHandler) CreateUserHandler(c echo.Context) error {
	user := new(service.User)

	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(
			echo.ErrInternalServerError.Code,
			fmt.Sprintf("Something went wrong while binding User: %v", err),
		)
	}

	if c.Request().Method == "POST" {
		if err := u.UserService.CreateUser(user); err != nil {
			return echo.NewHTTPError(
				echo.ErrInternalServerError.Code,
				fmt.Sprintf("Something went wrong while trying to create User: %e", err),
			)
		}
	}

	users, err := u.UserService.GetAllUsers()
	if err != nil {
		return echo.NewHTTPError(
			echo.ErrInternalServerError.Code,
			fmt.Sprintf("Something went wrong while trying to create User: %e", err),
		)
	}

	return web.UsersList(users).Render(c.Request().Context(), c.Response().Writer)
}

func (u *UserHandler) ShowUsers(c echo.Context) error {
	if err := web.UserForm().Render(c.Request().Context(), c.Response().Writer); err != nil {
		return echo.NewHTTPError(
			echo.ErrInternalServerError.Code,
			fmt.Sprintf("Error rendering users list: %e", err),
		)
	}

	return nil
}
