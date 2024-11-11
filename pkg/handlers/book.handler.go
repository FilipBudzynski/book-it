package handlers

import (
	"net/http"

	web_books "github.com/FilipBudzynski/book_it/cmd/web/books"
	"github.com/FilipBudzynski/book_it/pkg/models"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/labstack/echo/v4"
)

// BookService provides actions for managing book resources.
type BookService interface {
	// GetByTitle returns maxResults number of books by title from external api
	GetByTitle(title string, maxResults int) ([]*models.Book, error)
	// GetMaxResults gets the maxResults value specified for the service
	GetMaxResults() int
}

type BookHandler struct {
	bookService BookService
}

func NewBookHandler(bh BookService) *BookHandler {
	return &BookHandler{
		bookService: bh,
	}
}

func (h *BookHandler) Search(c echo.Context) error {
	if c.Request().Method == "GET" {
		return utils.RenderView(c, web_books.BooksSearch())
	}

	query := c.FormValue("book-title")

	exampleBooks, err := h.bookService.GetByTitle(query, h.bookService.GetMaxResults())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return utils.RenderView(c, web_books.BooksPost(exampleBooks))
}

func (h *BookHandler) List(c echo.Context) error {
	// bookId := c.Param("book_id")
	// limit, err := strconv.Atoi(c.Param("limit"))
	// if err != nil {
	// 	return echo.NewHTTPError(echo.ErrUnprocessableEntity.Code, err.Error())
	// }
	//
	// books, err := h.bookService.GetWithMaxResults(bookId, limit)
	//    if err != nil {
	//
	//    }
	// _ = books
	//
	// return utils.RenderView(c, web.ListBooks(books))
	return nil
}
