package handlers

import (
	"net/http"

	web_user_books "github.com/FilipBudzynski/book_it/cmd/web/user_books"
	"github.com/FilipBudzynski/book_it/pkg/models"
	"github.com/FilipBudzynski/book_it/pkg/schemas"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/labstack/echo/v4"
)

// UserBookService provides actions for managing user_book resources
type UserBookService interface {
	Create(userId, bookId string) error
	GetAll(userId string) ([]models.UserBook, error)
	GetById(id string) (*models.UserBook, error)
	Update(userBook *models.UserBook) error
	GetUserBooks(userId string) ([]schemas.Book, error)
}

type UserBookHandler struct {
	userBookService UserBookService
}

func NewUserBookHandler(userBookService UserBookService) *UserBookHandler {
	return &UserBookHandler{
		userBookService: userBookService,
	}
}

func (h *UserBookHandler) AddBook(c echo.Context) error {
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

func (h *UserBookHandler) List(c echo.Context) error {
	userId, err := utils.GetUserIDFromSession(c.Request())
	if err != nil {
		return echo.NewHTTPError(echo.ErrUnauthorized.Code, err.Error())
	}

	userBooks, err := h.userBookService.GetUserBooks(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return utils.RenderView(c, web_user_books.List(userBooks))
}
