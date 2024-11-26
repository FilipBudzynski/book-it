package handlers

import (
	"net/http"
	"net/url"

	web_books "github.com/FilipBudzynski/book_it/cmd/web/books"
	"github.com/FilipBudzynski/book_it/pkg/schemas"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/labstack/echo/v4"
)

// BookService provides actions for managing book resources.
//
// BookService can be implemented by different providers (e.g. Google Books API)
//
// It should communicate with the external api in order to retreive response
// and parse it into models.Book struct
type BookService interface {
	// GetByQuery returns maxResults number of books by title from external api
	GetByQuery(title string, maxResults int) ([]*schemas.Book, error)
	// GetMaxResults gets the maxResults value specified for the service
	GetMaxResults() int
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

	exampleBooks, err := h.bookService.GetByQuery(encodedQuery, h.bookService.GetMaxResults())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return utils.RenderView(c, web_books.BooksPost(exampleBooks))
}

func (h *BookHandler) List(c echo.Context) error {
	return nil
}
