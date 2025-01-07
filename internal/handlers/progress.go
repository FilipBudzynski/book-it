package handlers

import (
	"net/http"
	"strconv"
	"time"

	webProgress "github.com/FilipBudzynski/book_it/cmd/web/progress"
	"github.com/FilipBudzynski/book_it/internal/errs"
	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/internal/toast"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/labstack/echo/v4"
)

const (
	CompletedBookMessage  = "CONGRATULATIONS! You have completed the book!"
	TrackingBeginsMessage = "Tracking Begins!"
)

type ProgressService interface {
	// standard methods
	Create(bookId uint, totalPages int, bookTitle, startDateString, endDateString string) (models.ReadingProgress, error)
	Get(id string) (*models.ReadingProgress, error)
	GetByUserBookId(userBookId string) (*models.ReadingProgress, error)
	GetProgressAssosiatedWithLogId(id string) (*models.ReadingProgress, error)
	UpdateTargetPages(progressId uint, logDate time.Time) error
	Delete(id string) error

	// log methods
	GetLog(id string) (*models.DailyProgressLog, error)
	UpdateLogPagesRead(id string, pagesRead int) (*models.DailyProgressLog, error)
}

type progressHandler struct {
	progressService ProgressService
}

func NewProgressHandler(s ProgressService) *progressHandler {
	return &progressHandler{
		progressService: s,
	}
}

func (h *progressHandler) RegisterRoutes(app *echo.Echo) {
	group := app.Group("/progress")
	// middleware for protected routes
	group.Use(utils.CheckLoggedInMiddleware)
	// progress endpoints
	group.POST("", h.Create)
	group.GET("/:id", h.GetByUserBookId)
	// group.PUT("/:id", h.Edit)
	group.DELETE("/:id", h.Delete)

	// progress log endpoints
	group.Use(utils.CheckLoggedInMiddleware)
	group.POST("/log/:id", h.UpdatePagesRead)
	// htmx routes
	group.GET("/log/modal/:id", h.GetLogModal)
}

func (h *progressHandler) Create(c echo.Context) error {
	progressBind := &models.ReadingProgress{}
	if err := c.Bind(progressBind); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	startDateString := c.FormValue("start-date")
	endDateString := c.FormValue("end-date")

	progress, err := h.progressService.Create(
		progressBind.UserBookID,
		progressBind.TotalPages,
		progressBind.BookTitle,
		startDateString,
		endDateString,
	)
	if err != nil {
		return errs.HttpErrorBadRequest(err)
	}

	_ = toast.Success(c, TrackingBeginsMessage)
	return utils.RenderView(c, webProgress.TrackingButton(progress.UserBookID, progress.Completed))
}

func (h *progressHandler) GetByUserBookId(c echo.Context) error {
	id := c.Param("id")
	progress, err := h.progressService.GetByUserBookId(id)
	if err != nil {
		return errs.HttpErrorBadRequest(err)
	}
	return utils.RenderView(c, webProgress.ProgressStatistics(progress))
}

func (h *progressHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	err := h.progressService.Delete(id)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	c.Response().Header().Set("HX-Redirect", "/user-books")
	return c.NoContent(http.StatusOK)
}

func (h *progressHandler) UpdatePagesRead(c echo.Context) error {
	id := c.Param("id")
	pagesRead, err := strconv.Atoi(c.FormValue("pages-read"))
	if err != nil {
		return errs.HttpErrorBadRequest(err)
	}

	log, err := h.progressService.UpdateLogPagesRead(id, pagesRead)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	progress, err := h.progressService.GetProgressAssosiatedWithLogId(id)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	if !progress.IsFinishedOnLastLog(log.Date) {
		toast.Info(models.ErrProgressLastDayNotFinished.Error()).SetHXTriggerHeader(c)
	}

	if progress.Completed {
		toast.Success(c, CompletedBookMessage)
	}

	return utils.RenderView(c, webProgress.ProgressStatistics(progress))
}

func (h *progressHandler) GetLogModal(c echo.Context) error {
	id := c.Param("id")
	log, err := h.progressService.GetLog(id)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}
	return utils.RenderView(c, webProgress.ProgressLogModal(*log))
}
