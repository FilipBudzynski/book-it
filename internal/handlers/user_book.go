package handlers

import (
	"fmt"
	"net/http"

	webBooks "github.com/FilipBudzynski/book_it/cmd/web/books"
	webUserBooks "github.com/FilipBudzynski/book_it/cmd/web/user_books"
	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/labstack/echo/v4"
)

// UserBookService provides actions for managing user_book resources
type UserBookService interface {
	Create(userId, bookId string) error
	Update(userBook *models.UserBook) error
	Delete(id string) error
	GetAll(userId string) ([]*models.UserBook, error)
	GetById(id string) (*models.UserBook, error)
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

	return utils.RenderView(c, webBooks.WantToReadButton(bookID, true))
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

	return utils.RenderView(c, webBooks.WantToReadButton(bookID, false))
}

func (h *UserBookHandler) List(c echo.Context) error {
	userId, err := utils.GetUserIDFromSession(c.Request())
	if err != nil {
		return echo.NewHTTPError(echo.ErrUnauthorized.Code, err.Error())
	}

	userBooks, err := h.userBookService.GetAll(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return utils.RenderView(c, webUserBooks.List(userBooks))
}

func (h *UserBookHandler) GetCreateTrackingModal(c echo.Context) error {
	bookID := c.Param("user_book_id")
	if bookID == "" {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			fmt.Errorf("Something went wrong with the request. Book ID was not provided in query parameters"),
		)
	}

	userBook, err := h.userBookService.GetById(bookID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return utils.RenderView(c, webUserBooks.ProgressCreateModal(userBook))
}

func (h *UserBookHandler) RegisterRoutes(app *echo.Echo) {
	group := app.Group("/user-books")
	// middleware for protected routes
	group.Use(utils.CheckLoggedInMiddleware)
	// UserBook endpoints
	group.POST("/:book_id", h.Create)
	group.DELETE("/:book_id", h.Delete)
	group.GET("", h.List)
	group.GET("/create_modal/:user_book_id", h.GetCreateTrackingModal)
}
