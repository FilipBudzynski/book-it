package server

import (
	"net/http"

	"github.com/FilipBudzynski/book_it/cmd/web"
	"github.com/FilipBudzynski/book_it/pkg/handlers"
	"github.com/FilipBudzynski/book_it/pkg/routes"
	"github.com/FilipBudzynski/book_it/pkg/services"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

var UserService handlers.UserService

func (s *Server) RegisterRoutes(db *gorm.DB) http.Handler {
	e := echo.New()
	UserService = services.NewUserService(db)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	fileServer := http.FileServer(http.FS(web.Files))
	e.GET("/assets/*", echo.WrapHandler(fileServer))

	e.GET("/web", echo.WrapHandler(templ.Handler(web.HelloForm())))
	e.POST("/hello", echo.WrapHandler(http.HandlerFunc(web.HelloWebHandler)))

	// e.GET("/", s.HelloWorldHandler)
	e.GET("/", s.LandingPageHandler)
	e.GET("/health", s.healthHandler)

	// register user routes
	// userService := services.NewUserService(db)
	userHandler := handlers.NewUserHandler(UserService)
	routes.RegisterUserRoutes(e, userHandler)

	// register auth routes
	authHanlder := handlers.NewAuthHandler(UserService)
	routes.RegisterAuthRoutes(e, authHanlder)

	// register book routes
	googleBookService := services.NewGoogleBookService()
	bookHanlder := handlers.NewBookHandler(googleBookService)
	routes.RegisterBookRoutes(e, bookHanlder)

	return e
}

func (s *Server) LandingPageHandler(c echo.Context) error {
	userSession, _ := utils.GetUserSessionFromStore(c.Request())
	if (userSession == utils.UserSession{}) {
		return utils.RenderView(c, web.HomePage(nil))
	}

	dbUser, err := UserService.GetByGoogleID(userSession.UserID)
	if dbUser == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return utils.RenderView(c, web.HomePage(dbUser))
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}
