package handlers

import (
	"net/http"
	"net/url"

	web_books "github.com/FilipBudzynski/book_it/cmd/web/books"
	"github.com/FilipBudzynski/book_it/pkg/schemas"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/labstack/echo/v4"
)

const booksLimit int = 10

// BookService provides actions for managing book resources.
// BookSerice should uses a provider to get books from external APIs or database
type BookService interface {
	// GetByQuery returns maxResults number of books by title from external api
	GetByQuery(title string, maxResults int) ([]schemas.Book, error)
	// GetByID
	GetByID(id string) (schemas.Book, error)
}

type BookHandler struct {
	bookService BookService
}

func NewBookHandler(bh BookService) *BookHandler {
	return &BookHandler{
		bookService: bh,
	}
}

func (h *BookHandler) ListBooks(c echo.Context) error {
	if c.Request().Method == "GET" {
		return utils.RenderView(c, web_books.BooksSearch())
	}

	query := c.FormValue("book-title")
	encodedQuery := url.QueryEscape(query)

	books, err := h.bookService.GetByQuery(encodedQuery, booksLimit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return utils.RenderView(c, web_books.BooksPost(books))
}

func (h *BookHandler) List(c echo.Context) error {
	return nil
}
