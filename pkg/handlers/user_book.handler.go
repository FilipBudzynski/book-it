package handlers

import (
	"net/http"

	"github.com/FilipBudzynski/book_it/pkg/models"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/labstack/echo/v4"
)

// UserBookService provides actions for managing user_book resources
type UserBookService interface {
	Create(userId, bookId string) error
	GetAll() ([]models.UserBook, error)
	GetById(id string) (*models.UserBook, error)
	Update(userBook *models.UserBook) error
}

type UserBookHandler struct {
	userBookService UserBookService
}

func NewUserBookHandler(userBookService UserBookService) *UserBookHandler {
	return &UserBookHandler{
		userBookService: userBookService,
	}
}

func (h *UserBookHandler) Create(c echo.Context) error {
	bookId := c.QueryParam("book-id")
	userId, err := utils.GetUserIDFromSession(c.Request())
	if err != nil {
		return echo.NewHTTPError(echo.ErrUnauthorized.Code, err.Error())
	}

	err = h.userBookService.Create(userId, bookId)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}

	// TODO: render view
	return nil
}
