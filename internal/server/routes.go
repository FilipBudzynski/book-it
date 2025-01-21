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
	e.Use(utils.ErrorPagesMiddleware)
	e.Use(toast.ToastMiddleware)

	return s
}

var notifyManager *handlers.ConnectionManager

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

	notifyManager = handlers.NewConnectionManager()

	routeRegistrars := []RouteRegistrar{
		handlers.NewAuthHandler(userService),
		handlers.NewUserHandler(userService),
		handlers.NewBookHandler(bookService, userBookService, userService),
		handlers.NewUserBookHandler(userBookService),
		handlers.NewProgressHandler(progressService, userBookService),
		handlers.NewExchangeHandler(exchangeService, bookService, notifyManager),
	}

	for _, routeRegistrar := range routeRegistrars {
		routeRegistrar.RegisterRoutes(e)
	}

	fileServer := http.FileServer(http.FS(web.Files))
	e.GET("/assets/*", echo.WrapHandler(fileServer))
	e.Static("/static", "cmd/web")

	e.GET("/sse", notifyManager.SseHandler)

	return s
}

// func sseHandler(c echo.Context) error {
// 	// Set headers for SSE
// 	c.Response().Header().Set("Content-Type", "text/event-stream")
// 	c.Response().Header().Set("Cache-Control", "no-cache")
// 	c.Response().Header().Set("Connection", "keep-alive")
//
// 	// Create a channel to send data
// 	dataCh := make(chan string)
//
// 	// Create a context for handling client disconnection
// 	_, cancel := context.WithCancel(c.Request().Context())
// 	defer cancel()
//
// 	// Send data to the client
// 	go func() {
// 		for data := range dataCh {
// 			fmt.Fprintf(c.Response().Writer, "data: %s\n\n", data)
// 			c.Response().Writer.(http.Flusher).Flush()
// 		}
// 	}()
//
// 	// Simulate sending data periodically
// 	for {
// 		dataCh <- time.Now().Format(time.TimeOnly)
// 		time.Sleep(1 * time.Second)
// 	}
//
// }
