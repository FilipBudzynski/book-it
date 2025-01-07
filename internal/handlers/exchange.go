package handlers

import (
	"fmt"
	"net/http"
	"reflect"

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
	exchangeBind := &ExchangeFormBinding{}
	if err := exchangeBind.Bind(c); err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	userId, err := utils.GetUserIDFromSession(c.Request())
	if err != nil {
		return errs.HttpErrorUnauthorized(err)
	}

	if err := h.exchangeService.Create(userId, exchangeBind.DesiredBookID, exchangeBind.UserBookIDs); err != nil {
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

type ExchangeFormBinding struct {
	DesiredBookID string `form:"desired-book-id"`
	UserBookIDs   []string
}

func (e *ExchangeFormBinding) Bind(c echo.Context) error {
	val := reflect.ValueOf(e).Elem()
	typ := reflect.TypeOf(*e)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		formTag := field.Tag.Get("form")
		if formTag != "" {
			if value := c.FormValue(formTag); value != "" {
				val.Field(i).SetString(value)
			}
		}
	}
	e.UserBookIDs = []string{}
	for i := 0; i <= 4; i++ { // Adjust the range based on expected inputs
		fieldName := fmt.Sprintf("user-book-%d", i)
		if value := c.FormValue(fieldName); value != "" {
			e.UserBookIDs = append(e.UserBookIDs, value)
		}
	}
	return nil
}
