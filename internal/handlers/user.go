package handlers

import (
	"fmt"
	"net/http"
	"strconv"

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
	Delete(id string) error

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
	group.GET("/profile", h.Profile)
	group.GET("/profile/location/modal", h.GetLocationModal)
	group.POST("/profile/genres/:genre_id", h.AddGenre)
	group.DELETE("/profile/genres/:genre_id", h.RemoveGenre)
	group.DELETE("", h.Delete)
	group.POST("/profile/location", h.ChangeLocation)
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

func (h *UserHandler) GetLocationModal(c echo.Context) error {
	userID, err := utils.GetUserIDFromSession(c.Request())
	if err != nil {
		return errs.HttpErrorUnauthorized(err)
	}
	user, err := h.userService.GetByGoogleID(userID)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}
	return utils.RenderView(c, webUser.LocationModal(user))
}

func (h *UserHandler) AddGenre(c echo.Context) error {
	genreID := c.Param("genre_id")

	userID, err := utils.GetUserIDFromSession(c.Request())
	if err != nil {
		return errs.HttpErrorUnauthorized(err)
	}

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

func (h *UserHandler) LandingPage(c echo.Context) error {
	return utils.RenderView(c, web.HomePage())
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

func (h *UserHandler) Delete(c echo.Context) error {
	userID, err := utils.GetUserIDFromSession(c.Request())
	if err != nil {
		return errs.HttpErrorUnauthorized(err)
	}

	err = h.userService.Delete(userID)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	err = utils.RemoveCookieSession(c.Response().Writer, c.Request())
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusOK)
}

func (h *UserHandler) ChangeLocation(c echo.Context) error {
	lat := c.FormValue("latitude")
	lon := c.FormValue("longitude")
	formatted := c.FormValue("formatted")

	latParsed, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}
	lonParsed, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	userID, err := utils.GetUserIDFromSession(c.Request())
	if err != nil {
		return errs.HttpErrorUnauthorized(err)
	}
	user, err := h.userService.GetByGoogleID(userID)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	loc := &models.Location{
		Formatted: formatted,
		Latitude:  latParsed,
		Longitude: lonParsed,
	}
	user.Location = loc
	err = h.userService.Update(user)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	return c.NoContent(http.StatusOK)
}
