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

type ProgressService interface {
	// standard methods
	Create(progress *models.ReadingProgress) error
	Get(id string) (*models.ReadingProgress, error)
	Update(progress *models.ReadingProgress) error
	Delete(id string) error

	// custom methods
	GetByUserBookId(userBookId string) (*models.ReadingProgress, error)
}

type ProgressLogService interface {
	// standard methods
	Create(progressId, userBookId uint, target int, date time.Time) (*models.DailyProgressLog, error)
	Update(log *models.DailyProgressLog) error
	Get(id string) (*models.DailyProgressLog, error)
	Delete(id string) error
	GetAll(progressId string) ([]models.DailyProgressLog, error)
}

type progressHandler struct {
	progressService    ProgressService
	progressLogService ProgressLogService
}

func NewProgressHandler(s ProgressService) *progressHandler {
	return &progressHandler{
		progressService: s,
	}
}

func (h *progressHandler) WithProgressLogService(progressLogService ProgressLogService) *progressHandler {
	h.progressLogService = progressLogService
	return h
}

func (h *progressHandler) RegisterRoutes(app *echo.Echo) {
	group := app.Group("/progress")
	// middleware for protected routes
	group.Use(utils.CheckLoggedInMiddleware)
	// progress endpoints
	group.POST("", h.Create)
	group.GET("/:id", h.GetByUserBookId)
	group.PUT("", h.Update)
	group.DELETE("", h.Delete)

	// progress log endpoints
	group.Use(utils.CheckLoggedInMiddleware)
	group.POST("/log/:id", h.UpdateLog)
	// htmx routes
	group.GET("/log/modal/:id", h.GetLogModal)
	group.GET("/log/list/:id", h.GetLogList)
}

func (s *progressHandler) Create(c echo.Context) error {
	bookProgress := &models.ReadingProgress{}
	if err := c.Bind(bookProgress); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	startDate := c.FormValue("start-date")
	endDate := c.FormValue("end-date")

	startDateParsed, err := time.Parse(time.DateOnly, startDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	endDateParsed, err := time.Parse(time.DateOnly, endDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	days := int(endDateParsed.Sub(startDateParsed).Hours() / 24)
	if days == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "End date must be after start date")
	}
	bookProgress.DailyTargetPages = int(bookProgress.TotalPages / days)

	bookProgress.StartDate = startDateParsed
	bookProgress.EndDate = endDateParsed
	bookProgress.Completed = false


	for i := range days {
		trackingLog, err := s.progressLogService.Create(
			bookProgress.ID,
			bookProgress.UserBookID,
			bookProgress.DailyTargetPages,
			time.Now().AddDate(0, 0, i),
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		bookProgress.DailyProgress = append(bookProgress.DailyProgress, *trackingLog)
	}

	err = s.progressService.Create(bookProgress)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return utils.RenderView(c, webProgress.OnTrackIdentifiactor(bookProgress.UserBookID))
}

func (s *progressHandler) Update(c echo.Context) error {
	progress := &models.ReadingProgress{}
	if err := c.Bind(progress); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err := s.progressService.Update(progress)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, progress)
}

func (s *progressHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	err := s.progressService.Delete(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (s *progressHandler) GetByUserBookId(c echo.Context) error {
	id := c.Param("id")
	progress, err := s.progressService.GetByUserBookId(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return utils.RenderView(c, webProgress.ProgressStatistics(progress))
}

func (s *progressHandler) UpdateLog(c echo.Context) error {
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

	progress, err := s.progressService.Get(fmt.Sprintf("%d", dailyLog.ReadingProgressID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return utils.RenderView(c, webProgress.ProgressStatistics(progress))
}

// GetModal return reading log modal component if log with given id exists
func (s *progressHandler) GetLogModal(c echo.Context) error {
	id := c.Param("id")
	readingLog, err := s.progressLogService.Get(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return utils.RenderView(c, webProgress.ProgressLogModal(*readingLog))
}

func (s *progressHandler) GetLogList(c echo.Context) error {
	progressId := c.Param("id")
	dailyLogs, err := s.progressLogService.GetAll(progressId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return utils.RenderView(c, webProgress.DailyProgressLogs(dailyLogs))
}
