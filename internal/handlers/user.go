package handlers

import (
	"fmt"
	"net/http"

	"github.com/FilipBudzynski/book_it/cmd/web"
	webUser "github.com/FilipBudzynski/book_it/cmd/web/user"
	"github.com/FilipBudzynski/book_it/internal/errs"
	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/internal/toast"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/labstack/echo/v4"
)

// UserService provides actions for managing Users.
type UserService interface {
	Create(u *models.User) error
	Update(u *models.User) error
	GetById(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByGoogleID(googleID string) (*models.User, error)
	GetAll() ([]models.User, error)
	Delete(u models.User) error

	AddGenre(userID, genre string) (*models.Genre, error)
	RemoveGenre(userID, genre string) (*models.Genre, error)
	GetAllGenres() ([]*models.Genre, error)
}

type UserHandler struct {
	userService UserService
}

func NewUserHandler(us UserService) *UserHandler {
	return &UserHandler{
		userService: us,
	}
}

func (h *UserHandler) RegisterRoutes(app *echo.Echo) {
	app.GET("/", h.LandingPage)
	app.GET("/navbar", h.Navbar)
	group := app.Group("/users")
	group.Use(utils.CheckLoggedInMiddleware)
	group.GET("", h.ListUsers)
	group.GET("/profile", h.Profile)
	group.POST("/profile/genres/:genre_id", h.AddGenre)
	group.DELETE("/profile/genres/:genre_id", h.RemoveGenre)
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	user := new(models.User)

	if err := c.Bind(user); err != nil {
		return errs.HttpErrorBadRequest(err)
	}

	if err := h.userService.Create(user); err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	_ = toast.Success(c, "Account created")
	return c.NoContent(http.StatusCreated)
}

func (h *UserHandler) Profile(c echo.Context) error {
	userID, err := utils.GetUserIDFromSession(c.Request())
	if err != nil {
		return errs.HttpErrorUnauthorized(err)
	}

	user, err := h.userService.GetByGoogleID(userID)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	genres, err := h.userService.GetAllGenres()
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}
	for _, genre := range user.Genres {
		fmt.Printf("USER GENRE: %s\n", genre.Name)
	}

	return utils.RenderView(c, webUser.Profile(user, genres))
}

func (h *UserHandler) ListUsers(c echo.Context) error {
	users, err := h.userService.GetAll()
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	return utils.RenderView(c, web.UserForm(users))
}

func (h *UserHandler) AddGenre(c echo.Context) error {
	genreID := c.Param("genre_id")

	userID, err := utils.GetUserIDFromSession(c.Request())
	if err != nil {
		return errs.HttpErrorUnauthorized(err)
	}
	fmt.Println("USER ID: ", userID)

	genre, err := h.userService.AddGenre(userID, genreID)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	return utils.RenderView(c, webUser.GenreButton(genre, true))
}

func (h *UserHandler) RemoveGenre(c echo.Context) error {
	genreID := c.Param("genre_id")

	userID, err := utils.GetUserIDFromSession(c.Request())
	if err != nil {
		return errs.HttpErrorUnauthorized(err)
	}

	genre, err := h.userService.RemoveGenre(userID, genreID)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	return utils.RenderView(c, webUser.GenreButton(genre, false))
}

// LandingPageHandler returns a landing page
func (h *UserHandler) LandingPage(c echo.Context) error {
	userSession, _ := utils.GetUserSessionFromStore(c.Request())
	if (userSession == utils.UserSession{}) {
		return utils.RenderView(c, web.HomePage(nil))
	}

	dbUser, err := h.userService.GetByGoogleID(userSession.UserID)
	if dbUser == nil {
		return errs.HttpErrorInternalServerError(err)
	}

	return utils.RenderView(c, web.HomePage(dbUser))
}

func (h *UserHandler) Navbar(c echo.Context) error {
	userSession, err := utils.GetUserSessionFromStore(c.Request())
	if err != nil {
		return utils.RenderView(c, web.Navbar(nil))
	}

	user, err := h.userService.GetByGoogleID(userSession.UserID)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	return utils.RenderView(c, web.Navbar(user))
}
