package handlers

import (
	"fmt"

	"github.com/FilipBudzynski/book_it/cmd/web"
	"github.com/FilipBudzynski/book_it/pkg/models"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/labstack/echo/v4"
)

// UserService provides actions for managing Users.
type UserService interface {
	// db methods
	Create(u *models.User) error
	Update(u *models.User) error
	GetById(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByGoogleID(googleID string) (*models.User, error)
	GetAll() ([]models.User, error)
	Delete(u models.User) error
}

type UserHandler struct {
	userService UserService
}

func NewUserHandler(us UserService) *UserHandler {
	return &UserHandler{
		userService: us,
	}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	user := new(models.User)

	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, err)
	}

	if err := h.userService.Create(user); err != nil {
		return fmt.Errorf("creating user failed, err %v", err)
		// echo.NewHTTPError(echo.ErrInternalServerError.Code, err)
	}

	fmt.Fprintln(c.Response().Writer, "User registered successfully")
	return nil

	// return web.AppendUsersList(*user).Render(c.Request().Context(), c.Response().Writer)
}

func (h *UserHandler) ListUsers(c echo.Context) error {
	users, err := h.userService.GetAll()
	if err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, err)
	}

	return utils.RenderView(c, web.UserForm(users))
}

func (h *UserHandler) Navbar(c echo.Context) error {
	userSession, err := utils.GetUserSessionFromStore(c.Request())
	if err != nil {
		return utils.RenderView(c, web.Navbar(nil))
	}

	user, err := h.userService.GetByGoogleID(userSession.UserID)
	if err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, err)
	}

	return utils.RenderView(c, web.Navbar(user))
}
