package handlers

import (
	"net/http"

	webBooks "github.com/FilipBudzynski/book_it/cmd/web/books"
	webExchange "github.com/FilipBudzynski/book_it/cmd/web/exchange"
	webProgress "github.com/FilipBudzynski/book_it/cmd/web/progress"
	webUserBooks "github.com/FilipBudzynski/book_it/cmd/web/user_books"
	"github.com/FilipBudzynski/book_it/internal/errs"
	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/labstack/echo/v4"
)

// UserBookService provides actions for managing user_book resources
type UserBookService interface {
	Create(userId, bookId string) error
	Update(userBook *models.UserBook) error
	Get(id string) (*models.UserBook, error)
	GetAll(userId string) ([]*models.UserBook, error)
	Delete(id string) error
	DeleteByBookId(bookId string) error
	Search(userId , query string) ([]*models.UserBook, error)
}

type UserBookHandler struct {
	userBookService UserBookService
}

func NewUserBookHandler(userBookService UserBookService) *UserBookHandler {
	return &UserBookHandler{
		userBookService: userBookService,
	}
}

func (h *UserBookHandler) RegisterRoutes(app *echo.Echo) {
	group := app.Group("/user-books")
	group.Use(utils.CheckLoggedInMiddleware) // middleware for protected routes
	group.POST("/:book_id", h.Create)
	group.DELETE("/:book_id", h.Delete)
	group.DELETE("/search/:book_id", h.DeleteAndReplaceButton)
	group.GET("", h.List)
	group.GET("/create_modal/:user_book_id", h.GetCreateProgressModal)
	// group.GET("/exchange/additional", h.GetAdditionalOfferedBookInput)
	group.GET("/exchange/books", h.GetOfferedBooks)
	group.GET("/search", h.Search)
}

func (h *UserBookHandler) Create(c echo.Context) error {
	bookID := c.Param("book_id")
	if bookID == "" {
		return errs.HttpErrorBadRequest(models.ErrUserBookQueryWithoutId)
	}

	userID, err := utils.GetUserIDFromSession(c.Request())
	if err != nil {
		return errs.HttpErrorUnauthorized(err)
	}

	if err = h.userBookService.Create(userID, bookID); err != nil {
		return errs.HttpErrorConflict(err)
	}

	return utils.RenderView(c, webBooks.WantToReadButton(bookID, true))
}

func (h *UserBookHandler) Delete(c echo.Context) error {
	bookID := c.Param("book_id")
	if bookID == "" {
		return errs.HttpErrorBadRequest(models.ErrUserBookQueryWithoutId)
	}

	if err := h.userBookService.Delete(bookID); err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	return c.NoContent(http.StatusOK)
}

func (h *UserBookHandler) DeleteAndReplaceButton(c echo.Context) error {
	bookID := c.Param("book_id")
	if bookID == "" {
		return errs.HttpErrorBadRequest(models.ErrUserBookQueryWithoutId)
	}

	if err := h.userBookService.DeleteByBookId(bookID); err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	return utils.RenderView(c, webBooks.WantToReadButton(bookID, false))
}

func (h *UserBookHandler) List(c echo.Context) error {
	userId, err := utils.GetUserIDFromSession(c.Request())
	if err != nil {
		return errs.HttpErrorUnauthorized(err)
	}

	userBooks, err := h.userBookService.GetAll(userId)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	return utils.RenderView(c, webUserBooks.List(userBooks))
}

func (h *UserBookHandler) GetCreateProgressModal(c echo.Context) error {
	bookID := c.Param("user_book_id")
	if bookID == "" {
		return errs.HttpErrorBadRequest(models.ErrUserBookQueryWithoutId)
	}

	userBook, err := h.userBookService.Get(bookID)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	return utils.RenderView(c, webProgress.ProgressCreateModal(userBook))
}

func (h *UserBookHandler) GetOfferedBooks(c echo.Context) error {
	userId, err := utils.GetUserIDFromSession(c.Request())
	if err != nil {
		return errs.HttpErrorUnauthorized(err)
	}

	userBooks, err := h.userBookService.GetAll(userId)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	return utils.RenderView(c, webExchange.OfferedBooks(userBooks))
}

func (h *UserBookHandler) Search(c echo.Context) error {
	userId, err := utils.GetUserIDFromSession(c.Request())
	if err != nil {
		return errs.HttpErrorUnauthorized(err)
	}
	search := c.QueryParam("query")

    results, err := h.userBookService.Search(userId, search)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}
	return utils.RenderView(c, webUserBooks.BooksTableRows(results))
}
