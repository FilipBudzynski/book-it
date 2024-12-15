package routes

import (
	"github.com/FilipBudzynski/book_it/internal/handlers"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/labstack/echo/v4"
)

func RegisterTrackingRoutes(app *echo.Echo, h *handlers.TrackingHandler) {
	group := app.Group("/tracking")
	// middleware for protected routes
	group.Use(utils.CheckLoggedInMiddleware)
	// UserBook endpoints
	group.POST("", h.Create)
	group.GET("/:id", h.GetByUserBookId)
	group.PUT("", h.Update)
	group.DELETE("", h.Delete)
}
