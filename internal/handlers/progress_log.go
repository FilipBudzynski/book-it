package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	webProgress "github.com/FilipBudzynski/book_it/cmd/web/progress"
	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/labstack/echo/v4"
)

type ProgressLogService interface {
	// standard methods
	Create(progressId uint, target int, date time.Time) (*models.DailyProgressLog, error)
	Update(log *models.DailyProgressLog) error
	Get(id string) (*models.DailyProgressLog, error)
	Delete(id string) error
	GetAll(progressId string) ([]models.DailyProgressLog, error)
}

type progressLogHandler struct {
	progressLogService ProgressLogService
}

func NewProgressLogHandler(readingLogService ProgressLogService) *progressLogHandler {
	return &progressLogHandler{
		progressLogService: readingLogService,
	}
}

func (h *progressLogHandler) RegisterRoutes(app *echo.Echo) {
	group := app.Group("/progress_log")
	// middleware for protected routes
	group.Use(utils.CheckLoggedInMiddleware)
	group.PUT("/:id", h.Update)
	// htmx routes
	group.GET("/modal/:id", h.GetModal)
	group.GET("/list/:id", h.GetList)
}

func (s *progressLogHandler) Update(c echo.Context) error {
	id := c.Param("id")
	dailyLog, err := s.progressLogService.Get(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	pagesRead := c.FormValue("pages-read")
	if pagesRead == "" {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Something went wrong with the request. Pages read was not provided in form"))
	}

	pagesReadInt, err := strconv.Atoi(pagesRead)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if pagesReadInt < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Pages read cannot be negative")
	}
    // TODO: add check with total pages
	// if pagesReadInt > dailyLog.TargetPages {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "Pages read cannot be greater than target pages")
	// }

	dailyLog.PagesRead = int(pagesReadInt)

	if err = s.progressLogService.Update(dailyLog); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return utils.RenderView(c, webProgress.ProgressStep(*dailyLog))
}

// GetModal return reading log modal component if log with given id exists
func (s *progressLogHandler) GetModal(c echo.Context) error {
	id := c.Param("id")
	readingLog, err := s.progressLogService.Get(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return utils.RenderView(c, webProgress.ProgressLogModal(*readingLog))
}

func (s *progressLogHandler) GetList(c echo.Context) error {
	progressId := c.Param("id")
	dailyLogs, err := s.progressLogService.GetAll(progressId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return utils.RenderView(c, webProgress.DailyProgressLogs(dailyLogs))
}
