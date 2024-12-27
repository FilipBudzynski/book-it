package server

import (
	"net/http"

	"github.com/FilipBudzynski/book_it/cmd/web"
	"github.com/FilipBudzynski/book_it/internal/handlers"
	"github.com/FilipBudzynski/book_it/internal/providers"
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
	return s
}

func (s *Server) WithRegisterRoutes(e *echo.Echo) *Server {
	db := s.db.Db
	userService := services.NewUserService(db)
	bookService := services.NewBookService(db).
		WithProvider(providers.NewGoogleProvider().
			WithLimit(15))
	userBookService := services.NewUserBookService(db)
	progressService := services.NewProgressService(db)
	progressLogService := services.NewProgressLogService(db)

	routeRegistrars := []RouteRegistrar{
		handlers.NewAuthHandler(userService),
		handlers.NewUserHandler(userService),
		handlers.NewBookHandler(bookService, userBookService),
		handlers.NewUserBookHandler(userBookService),
		handlers.NewProgressHandler(progressService).WithProgressLogService(progressLogService),
		handlers.NewProgressLogHandler(progressLogService),
	}

	for _, routeRegistrar := range routeRegistrars {
		routeRegistrar.RegisterRoutes(e)
	}

	fileServer := http.FileServer(http.FS(web.Files))
	e.GET("/assets/*", echo.WrapHandler(fileServer))

	return s
}
