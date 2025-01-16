package server

import (
	"net/http"

	"github.com/FilipBudzynski/book_it/cmd/web"
	"github.com/FilipBudzynski/book_it/internal/handlers"
	"github.com/FilipBudzynski/book_it/internal/providers"
	"github.com/FilipBudzynski/book_it/internal/repositories"
	"github.com/FilipBudzynski/book_it/internal/services"
	"github.com/FilipBudzynski/book_it/internal/toast"
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
	e.Use(toast.ToastMiddleware)
	// e.Use(utils.ErrorPagesMiddleware)
	//
	// cors := cors.New(cors.Options{
	//         AllowedOrigins:   []string{"http://localhost:3000"},
	//         AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	//         AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
	//         AllowCredentials: true,
	//         MaxAge:           300, // Maximum age for preflight requests
	//     })
	// })
	//
	// e.Use(cors.Hanlder)

	return s
}

func (s *Server) WithRegisterRoutes(e *echo.Echo) *Server {
	db := s.db.Db

	progressRepo := repositories.NewProgressRepository(db)
	userRepo := repositories.NewUserRepository(db)
	userBookRepo := repositories.NewUserBookRepository(db)
	exchangeRequestRepo := repositories.NewExchangeRequestRepository(db)

	userService := services.NewUserService(userRepo)
	userBookService := services.NewUserBookService(userBookRepo, exchangeRequestRepo)
	progressService := services.NewProgressService(progressRepo)
	bookService := services.NewBookService(db).
		WithProvider(providers.NewGoogleProvider())
	exchangeService := services.NewExchangeService(exchangeRequestRepo)

	routeRegistrars := []RouteRegistrar{
		handlers.NewAuthHandler(userService),
		handlers.NewUserHandler(userService),
		handlers.NewBookHandler(bookService, userBookService, userService),
		handlers.NewUserBookHandler(userBookService),
		handlers.NewProgressHandler(progressService),
		handlers.NewExchangeHandler(exchangeService),
	}

	for _, routeRegistrar := range routeRegistrars {
		routeRegistrar.RegisterRoutes(e)
	}

	fileServer := http.FileServer(http.FS(web.Files))
	e.GET("/assets/*", echo.WrapHandler(fileServer))
	e.Static("/static", "cmd/web")

	return s
}
