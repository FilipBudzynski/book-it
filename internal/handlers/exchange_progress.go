package handlers

import (
	"net/http"

	webExchange "github.com/FilipBudzynski/book_it/cmd/web/exchange"
	"github.com/FilipBudzynski/book_it/internal/errs"
	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/internal/toast"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/labstack/echo/v4"
)

type ExchangeService interface {
	Create(userId, desiredBookID string, userBookIDs []string) (*models.ExchangeRequest, error)
	Get(id string) (*models.ExchangeRequest, error)
	GetAll(userId string) ([]*models.ExchangeRequest, error)
	Delete(id string) error
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
	group.GET("", h.Landing)
	group.GET("/modal/new", h.GetNewExchangeModal)
	group.GET("/list", h.GetAll)
	group.GET("/details/:id", h.Details)
	group.DELETE("/:id", h.Delete)
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

	exchange_request, err := h.exchangeService.Create(
		userId,
		exchangeBind.DesiredBookID,
		exchangeBind.UserBookIDs,
	)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	_ = toast.Success(c, "Exchange created!")
	return utils.RenderView(c, webExchange.ExchangeTableRow(*exchange_request))
}

func (h *exchangeHandler) Landing(c echo.Context) error {
	return utils.RenderView(c, webExchange.Landing())
}

func (h *exchangeHandler) GetAll(c echo.Context) error {
	userId, err := utils.GetUserIDFromSession(c.Request())
	if err != nil {
		return errs.HttpErrorUnauthorized(err)
	}

	exchanges, err := h.exchangeService.GetAll(userId)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	return utils.RenderView(c, webExchange.List(exchanges))
}

func (h *exchangeHandler) Details(c echo.Context) error {
	id := c.Param("id")
	exchange, err := h.exchangeService.Get(id)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}
	return utils.RenderView(c, webExchange.ExchangeDetails(exchange))
}

func (h *exchangeHandler) GetNewExchangeModal(c echo.Context) error {
	return utils.RenderView(c, webExchange.ExchangeModal())
}

func (h *exchangeHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	err := h.exchangeService.Delete(id)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}
	return c.NoContent(http.StatusNoContent)
}
