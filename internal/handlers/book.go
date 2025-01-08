package handlers

import (
	"net/http"
	"net/url"

	web_books "github.com/FilipBudzynski/book_it/cmd/web/books"
	"github.com/FilipBudzynski/book_it/internal/errs"
	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/labstack/echo/v4"
)

const booksLimit int = 10

// BookService provides actions for managing book resources.
// BookSerice should uses a provider to get books from external APIs or database
type BookService interface {
	Create(book *models.Book) error
	Delete(userID, bookID string) error
	// GetByQuery returns maxResults number of books by title from external api
	GetByQuery(title string, maxResults int) ([]*models.Book, error)
	// GetByID
	GetByID(id string) (*models.Book, error)

	WithProvider(provider BookProvider) BookService
}

// BookProvider is used to communicate with the external API or Database
// in order to retreive response and parse it into models.Book struct
type BookProvider interface {
	GetBook(id string) (*models.Book, error)
	GetBooksByQuery(query string, limit int) ([]*models.Book, error)
	// used to change the limit of query results
	WithLimit(limit int) BookProvider
}

type BookHandler struct {
	bookService      BookService
	userBooksService UserBookService
}

func NewBookHandler(bookService BookService, userBookService UserBookService) *BookHandler {
	return &BookHandler{
		bookService:      bookService,
		userBooksService: userBookService,
	}
}

func (h *BookHandler) RegisterRoutes(app *echo.Echo) {
	group := app.Group("/books")
	group.GET("", h.ListBooks)
	group.POST("", h.ListBooks)
	group.GET("/reduced/search", h.ReducedSearch)
	//group.POST("/reduced/search", h.ReducedSearch)
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

	userID, err := utils.GetUserIDFromSession(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	userBooks, err := h.userBooksService.GetAll(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return utils.RenderView(c, web_books.BooksPost(books, userBooks))
}

func (h *BookHandler) List(c echo.Context) error {
	return nil
}

func (h *BookHandler) ReducedSearch(c echo.Context) error {
	query := c.FormValue("book-title")
	encodedQuery := url.QueryEscape(query)

	books, err := h.bookService.GetByQuery(encodedQuery, booksLimit)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	return utils.RenderView(c, web_books.ReducedList(books))
}
