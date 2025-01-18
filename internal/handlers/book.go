package handlers

import (
	"net/http"
	"strconv"

	web_books "github.com/FilipBudzynski/book_it/cmd/web/books"
	webUser "github.com/FilipBudzynski/book_it/cmd/web/user"
	"github.com/FilipBudzynski/book_it/internal/errs"
	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/labstack/echo/v4"
)

// BookService provides actions for managing book resources.
// BookSerice should uses a provider to get books from external APIs or database
type BookService interface {
	Create(book *models.Book) error
	Delete(userID, bookID string) error
	// GetByQuery returns maxResults number of books by title from external api
	GetByQuery(query string, queryType QueryType, page int) ([]*models.Book, error)
	// GetByID
	GetByID(id string) (*models.Book, error)
	FetchReccomendations(genres []models.Genre) ([]*models.Book, error)

	WithProvider(provider BookProvider) BookService
	Provider() BookProvider
}

// BookProvider is used to communicate with the external API or Database
// in order to retreive response and parse it into models.Book struct
type BookProvider interface {
	GetBook(id string) (*models.Book, error)
	GetBooksByQuery(query string, queryType QueryType, limit, page int) ([]*models.Book, error)
	// used to change the limit of query results
	WithLimit(limit int) BookProvider
	GetLimit() int
	GetTotalForQuery(query string) int
	GetBooksByGenre(genre string, maxResults int) ([]*models.Book, error)
}

type BookHandler struct {
	bookService      BookService
	userService      UserService
	userBooksService UserBookService
}

func NewBookHandler(bookService BookService, userBookService UserBookService, userService UserService) *BookHandler {
	return &BookHandler{
		bookService:      bookService,
		userBooksService: userBookService,
		userService:      userService,
	}
}

func (h *BookHandler) RegisterRoutes(app *echo.Echo) {
	group := app.Group("/books")
	group.GET("", h.ListBooks)
	group.POST("", h.ListBooks)
	group.GET("/reduced/search", h.ReducedSearch)
	group.GET("/partial", h.BooksPartial)
	group.GET("/recommendations", h.Recommend)
}

func (h *BookHandler) ListBooks(c echo.Context) error {
	if c.Request().Method == "GET" {
		return utils.RenderView(c, web_books.BooksSearch())
	}

	query := c.FormValue("query")
	queryType := c.FormValue("type")

	books, userBooks, err := h.getBooksAndUserData(c, query, queryType, 1)
	if err != nil {
		return err
	}
	return utils.RenderView(c, web_books.BooksPost(books, userBooks, 2, query))
}

func (h *BookHandler) BooksPartial(c echo.Context) error {
	query := c.QueryParam("query")
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	books, userBooks, err := h.getBooksAndUserData(c, query, "title", page)
	if err != nil {
		return err
	}
	if len(books) == 0 {
		return c.NoContent(http.StatusNoContent)
	}

	return utils.RenderView(c, web_books.BooksPost(books, userBooks, page+1, query))
}

func (h *BookHandler) List(c echo.Context) error {
	return nil
}

func (h *BookHandler) ReducedSearch(c echo.Context) error {
	query := c.FormValue("book-title")
	books, err := h.bookService.GetByQuery(query, QueryTypeTitle, 1)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	return utils.RenderView(c, web_books.ReducedList(books))
}

func (g *BookHandler) Recommend(c echo.Context) error {
	userID, err := utils.GetUserIDFromSession(c.Request())
	if err != nil {
		return errs.HttpErrorUnauthorized(err)
	}

	user, err := g.userService.GetByGoogleID(userID)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	if len(user.Genres) == 0 {
		return utils.RenderView(c, webUser.Recommendations(nil))
	}
	recommendedBooks, err := g.bookService.FetchReccomendations(user.Genres)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	return utils.RenderView(c, webUser.Recommendations(recommendedBooks))
}

func (h *BookHandler) getBooksAndUserData(c echo.Context, query, queryTypeString string, page int) ([]*models.Book, []*models.UserBook, error) {
	// encodedQuery := url.QueryEscape(query)
	queryType := stringToQueryType(queryTypeString)

	books, err := h.bookService.GetByQuery(query, queryType, page)
	if err != nil {
		return nil, nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	userID, err := utils.GetUserIDFromSession(c.Request())
	if err != nil {
		return nil, nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	userBooks, err := h.userBooksService.GetAll(userID)
	if err != nil {
		return nil, nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return books, userBooks, nil
}
