package handlers

import (
	"net/http"

	webExchange "github.com/FilipBudzynski/book_it/cmd/web/exchange"
	"github.com/FilipBudzynski/book_it/internal/errs"
	"github.com/FilipBudzynski/book_it/internal/toast"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/labstack/echo/v4"
)

type ExchangeService interface {
	Create(userId, desiredBookID string, userBookIDs []string) error
}

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
	group.POST("", h.CreateExchange)
	group.GET("", h.List)
	group.GET("/modal/new", h.GetNewExchangeModal)
}

func (h *exchangeHandler) CreateExchange(c echo.Context) error {
	exchangeBind := &exchangeFormBinding{}
	if err := exchangeBind.bind(c); err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	userId, err := utils.GetUserIDFromSession(c.Request())
	if err != nil {
		return errs.HttpErrorUnauthorized(err)
	}

	err = h.exchangeService.Create(
		userId,
		exchangeBind.DesiredBookID,
		exchangeBind.UserBookIDs,
	)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	_ = toast.Success(c, "Exchange created!")
	return c.NoContent(http.StatusCreated)
}

func (h *exchangeHandler) List(c echo.Context) error {
	return utils.RenderView(c, webExchange.List())
}

func (h *exchangeHandler) GetNewExchangeModal(c echo.Context) error {
	return utils.RenderView(c, webExchange.ExchangeModal())
}
