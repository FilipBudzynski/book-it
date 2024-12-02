package handlers

import (
	"fmt"
	"net/http"

	web_books "github.com/FilipBudzynski/book_it/cmd/web/books"
	web_user_books "github.com/FilipBudzynski/book_it/cmd/web/user_books"
	"github.com/FilipBudzynski/book_it/pkg/models"
	"github.com/FilipBudzynski/book_it/pkg/schemas"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/labstack/echo/v4"
)

// UserBookService provides actions for managing user_book resources
type UserBookService interface {
	Create(userId, bookId string) error
	Update(userBook *models.UserBook) error
	Delete(id string) error
	GetAll(userId string) ([]models.UserBook, error)
	GetById(id string) (*models.UserBook, error)
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

func (h *UserBookHandler) Create(c echo.Context) error {
	bookID := c.Param("book_id")
	if bookID == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Something went wrong with the request. Book ID was not provided in query parameters"))
	}

	userID, err := utils.GetUserIDFromSession(c.Request())
	if err != nil {
		return echo.NewHTTPError(echo.ErrUnauthorized.Code, err.Error())
	}

	err = h.userBookService.Create(userID, bookID)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}

	return utils.RenderView(c, web_books.WantToReadButton(bookID, true))
}

func (h *UserBookHandler) Delete(c echo.Context) error {
	bookID := c.Param("book_id")
	if bookID == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Something went wrong with the request. Book ID was not provided in query parameters"))
	}

	err := h.userBookService.Delete(bookID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return utils.RenderView(c, web_books.WantToReadButton(bookID, false))
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
