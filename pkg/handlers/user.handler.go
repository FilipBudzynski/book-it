package handlers

import (
	"fmt"
	"net/http"

	"github.com/FilipBudzynski/book_it/cmd/web"
	"github.com/FilipBudzynski/book_it/pkg/models"
	"github.com/FilipBudzynski/book_it/pkg/services"
	"github.com/FilipBudzynski/book_it/utils"
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

func (h *UserHandler) AddBook(c echo.Context) error {
	bookId := c.QueryParam("book-id")
	userSession, err := utils.GetUserSessionFromStore(c.Request())
	if err != nil {
		return echo.NewHTTPError(echo.ErrUnauthorized.Code, err.Error())
	}

	user, err := h.userService.GetByGoogleID(userSession.UserID)
	if err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, err.Error())
	}

	err = h.userService.AddBook(user.GoogleId, bookId)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}

	return nil
}
