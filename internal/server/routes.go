package server

import (
	"book_it/cmd/web"
	"book_it/pkg/handlers"
	"book_it/pkg/services"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func (s *Server) RegisterRoutes(db *gorm.DB) http.Handler {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	fileServer := http.FileServer(http.FS(web.Files))
	e.GET("/assets/*", echo.WrapHandler(fileServer))

	e.GET("/web", echo.WrapHandler(templ.Handler(web.HelloForm())))
	e.POST("/hello", echo.WrapHandler(http.HandlerFunc(web.HelloWebHandler)))

	e.GET("/", s.HelloWorldHandler)

	e.GET("/health", s.healthHandler)

	userService := services.NewUserService(db)
	userHandler := handlers.NewUserHandler(userService)

    e.GET("/create", userHandler.ListUsersHandler)
	e.POST("/create", userHandler.CreateUserHandler)

	return e
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
