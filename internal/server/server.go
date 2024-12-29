package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/FilipBudzynski/book_it/internal/database"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
)

type Server struct {
	port int
	db   *database.Repository
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,
		db:   database.New(),
	}

	e := echo.New()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.WithMiddleware(e).WithRegisterRoutes(e).ToEchoHttpHandler(e),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

func (s *Server) ToEchoHttpHandler(e *echo.Echo) http.Handler {
	return e
}
