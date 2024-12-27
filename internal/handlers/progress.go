package handlers

import (
	"net/http"
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

type ProgressHandler struct {
	ProgressService    ProgressService
	ProgressLogService ProgressLogService
}

func NewProgressHandler(trackingService ProgressService) *ProgressHandler {
	return &ProgressHandler{
		ProgressService: trackingService,
	}
}

func (h *ProgressHandler) WithProgressLogService(progressLogService ProgressLogService) *ProgressHandler {
	h.ProgressLogService = progressLogService
	return h
}

func (s *ProgressHandler) Create(c echo.Context) error {
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
		trackingLog, err := s.ProgressLogService.Create(bookProgress.ID, bookProgress.DailyTargetPages, time.Now().AddDate(0, 0, i))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		bookProgress.DailyProgress = append(bookProgress.DailyProgress, *trackingLog)
	}

	err = s.ProgressService.Create(bookProgress)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return utils.RenderView(c, webProgress.OnTrackIdentifiactor(bookProgress.UserBookID))
}

func (s *ProgressHandler) Update(c echo.Context) error {
	progress := &models.ReadingProgress{}
	if err := c.Bind(progress); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err := s.ProgressService.Update(progress)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, progress)
}

func (s *ProgressHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	err := s.ProgressService.Delete(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (s *ProgressHandler) GetByUserBookId(c echo.Context) error {
	id := c.Param("id")
	progress, err := s.ProgressService.GetByUserBookId(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return utils.RenderView(c, webProgress.ProgressStatistics(progress))
}

// TODO: udpate daily log
func (s *ProgressHandler) UpdateDailyLog(c echo.Context) error {
	dailyLog := &models.DailyProgressLog{}
	if err := c.Bind(dailyLog); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return nil
}

func (h *ProgressHandler) RegisterRoutes(app *echo.Echo) {
	group := app.Group("/tracking")
	// middleware for protected routes
	group.Use(utils.CheckLoggedInMiddleware)

	// progress endpoints
	group.POST("", h.Create)
	group.GET("/:id", h.GetByUserBookId)
	group.PUT("", h.Update)
	group.DELETE("", h.Delete)
}
