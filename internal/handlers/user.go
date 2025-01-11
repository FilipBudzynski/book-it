package handlers

import (
	"net/http"

	"github.com/FilipBudzynski/book_it/cmd/web"
	"github.com/FilipBudzynski/book_it/internal/errs"
	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/internal/toast"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/labstack/echo/v4"
)

// UserService provides actions for managing Users.
type UserService interface {
	// db methods
	Create(u *models.User) error
	Update(u *models.User) error
	GetById(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByGoogleID(googleID string) (*models.User, error)
	GetAll() ([]models.User, error)
	Delete(u models.User) error
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

func (h *UserHandler) ListUsers(c echo.Context) error {
	users, err := h.userService.GetAll()
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	return utils.RenderView(c, web.UserForm(users))
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
