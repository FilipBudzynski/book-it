package handlers

import (
	webExchange "github.com/FilipBudzynski/book_it/cmd/web/exchange"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/labstack/echo/v4"
)

type ExchangeService interface{}

type exchangeHandler struct {
	exchangeService ExchangeService
}

func NewExchangeHandler(exchangeService ExchangeService) *exchangeHandler {
	return &exchangeHandler{
		exchangeService: exchangeService,
	}
}

func (h *exchangeHandler) RegisterRoutes(app *echo.Echo) {
	group := app.Group("/exchange")
	group.Use(utils.CheckLoggedInMiddleware) // protected routes
	group.GET("", h.List)
}

func (h *exchangeHandler) List(c echo.Context) error {
	return utils.RenderView(c, webExchange.List())
}
