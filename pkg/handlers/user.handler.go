package handlers

import (
	"book_it/cmd/web"
	"book_it/pkg/services"

	"github.com/labstack/echo/v4"
)

type UserService interface {
	Create(u *services.User) error
	Update(u *services.User) error
	GetById(id uint) (*services.User, error)
	GetAll() ([]services.User, error)
	Delete(u services.User) error
}

func NewUserHandler(us UserService) *UserHandler {
	return &UserHandler{
		userService: us,
	}
}

type UserHandler struct {
	userService UserService
}

func (u *UserHandler) CreateUserHandler(c echo.Context) error {
	user := new(services.User)

	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, err)
	}

	if err := u.userService.Create(user); err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, err)
	}

	return web.AppendUsersList(*user).Render(c.Request().Context(), c.Response().Writer)
}

func (u *UserHandler) ListUsersHandler(c echo.Context) error {
	users, err := u.userService.GetAll()
	if err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, err)
	}

	return web.UserForm(users).Render(c.Request().Context(), c.Response().Writer)
}
