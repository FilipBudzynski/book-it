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

func (s *Server) RegisterRoutes(db *gorm.DB) http.Handler {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// e.Use(session.Middleware(store))

	fileServer := http.FileServer(http.FS(web.Files))
	e.GET("/assets/*", echo.WrapHandler(fileServer))

	e.GET("/web", echo.WrapHandler(templ.Handler(web.HelloForm())))
	e.POST("/hello", echo.WrapHandler(http.HandlerFunc(web.HelloWebHandler)))

	// e.GET("/", s.HelloWorldHandler)
	e.GET("/", s.HomePageHandler)
	e.GET("/health", s.healthHandler)

	// register user routes
	userService := services.NewUserService(db)
	userHandler := handlers.NewUserHandler(userService)
	routes.RegisterUserRoutes(e, userHandler)

	// register auth routes
	authHanlder := handlers.NewAuthHandler(userService)
	routes.RegisterAuthRoutes(e, authHanlder)

	return e
}

func (s *Server) HomePageHandler(c echo.Context) error {
	logged := utils.IsUserLoggedIn(c.Request(), "session")
	return utils.RenderView(c, web.HomePage(logged))
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
