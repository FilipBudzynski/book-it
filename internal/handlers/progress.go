package handlers

import (
	"errors"
	"net/http"

	webProgress "github.com/FilipBudzynski/book_it/cmd/web/progress"
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
	Delete(id string) error

	// log methods
	GetLog(id string) (*models.DailyProgressLog, error)
	UpdateLogPagesRead(id, pagesReadString string) error
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
	group.DELETE("/:id", h.Delete)

	// progress log endpoints
	group.Use(utils.CheckLoggedInMiddleware)
	group.POST("/log/:id", h.UpdatePagesRead)
	// htmx routes
	group.GET("/log/modal/:id", h.GetLogModal)
}

func (s *progressHandler) Create(c echo.Context) error {
	progressBind := &models.ReadingProgress{}
	if err := c.Bind(progressBind); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	startDateString := c.FormValue("start-date")
	endDateString := c.FormValue("end-date")

	progress, err := s.progressService.Create(
		progressBind.UserBookID,
		progressBind.TotalPages,
		progressBind.BookTitle,
		startDateString,
		endDateString,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, toast.Warning(c, err.Error()))
	}

	toast.Success(c, TrackingBeginsMessage)
	return utils.RenderView(c, webProgress.TrackingButton(progress.UserBookID))
}

func (s *progressHandler) GetByUserBookId(c echo.Context) error {
	id := c.Param("id")
	progress, err := s.progressService.GetByUserBookId(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, toast.Warning(c, err.Error()))
	}
	return utils.RenderView(c, webProgress.ProgressStatistics(progress))
}

func (s *progressHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	err := s.progressService.Delete(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (s *progressHandler) UpdatePagesRead(c echo.Context) error {
	id := c.Param("id")
	pagesRead := c.FormValue("pages-read")

	err := s.progressService.UpdateLogPagesRead(id, pagesRead)
	if errors.Is(err, models.ErrProgressPastEndDate) {
		_ = toast.Info(c, err.Error())
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, toast.Danger(c, err.Error()))
	}

	progress, err := s.progressService.GetProgressAssosiatedWithLogId(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if progress.Completed {
		toast.Success(c, CompletedBookMessage)
	}
	return utils.RenderView(c, webProgress.ProgressStatistics(progress))
}

func (s *progressHandler) GetLogModal(c echo.Context) error {
	id := c.Param("id")
	log, err := s.progressService.GetLog(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return utils.RenderView(c, webProgress.ProgressLogModal(*log))
}
