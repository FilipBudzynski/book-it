package server

import (
	"net/http"

	"github.com/FilipBudzynski/book_it/cmd/web"
	"github.com/FilipBudzynski/book_it/internal/handlers"
	"github.com/FilipBudzynski/book_it/internal/providers"
	"github.com/FilipBudzynski/book_it/internal/repositories"
	"github.com/FilipBudzynski/book_it/internal/services"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/labstack/echo/v4"
	prettylogger "github.com/rdbell/echo-pretty-logger"
)

type RouteRegistrar interface {
	RegisterRoutes(e *echo.Echo)
}

func (s *Server) WithMiddleware(e *echo.Echo) *Server {
	e.Use(prettylogger.Logger)
	e.Use(utils.CustomRecoverMiddleware)
	e.Use(utils.RefreshSessionMiddleware)
	// e.HTTPErrorHandler = utils.CustomErrorHandler
	return s
}

func (s *Server) WithRegisterRoutes(e *echo.Echo) *Server {
	db := s.db.Db

	progressRepo := repositories.NewProgressRepository(db)
	userRepo := repositories.NewUserRepository(db)

	userService := services.NewUserService(userRepo)
	userBookService := services.NewUserBookService(db)
	progressService := services.NewProgressService(progressRepo)
	bookService := services.NewBookService(db).
		WithProvider(providers.NewGoogleProvider().
			WithLimit(15))

	routeRegistrars := []RouteRegistrar{
		handlers.NewAuthHandler(userService),
		handlers.NewUserHandler(userService),
		handlers.NewBookHandler(bookService, userBookService),
		handlers.NewUserBookHandler(userBookService),
		handlers.NewProgressHandler(progressService),
	}

	for _, routeRegistrar := range routeRegistrars {
		routeRegistrar.RegisterRoutes(e)
	}

	fileServer := http.FileServer(http.FS(web.Files))
	e.GET("/assets/*", echo.WrapHandler(fileServer))
	e.Static("/static", "cmd/web")

	return s
}
