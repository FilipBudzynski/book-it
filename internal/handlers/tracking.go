package handlers

import (
	"net/http"
	"time"

	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/labstack/echo/v4"
)

type TrackingService interface {
	// standard methods
	Create(progress *models.ReadingProgress) error
	Read(id string) (*models.ReadingProgress, error)
	Update(progress *models.ReadingProgress) error
	Delete(id string) error

	// custom methods
	GetByUserBookId(id string) (*models.ReadingProgress, error)
	GetDailyLog(bookId string, date time.Time) (*models.DailyReadingLog, error)
}

type TrackingHandler struct {
	TrackingService TrackingService
}

func NewTrackingHandler(trackingService TrackingService) TrackingHandler {
	return TrackingHandler{
		TrackingService: trackingService,
	}
}

func (s *TrackingHandler) GetDailyLog(c echo.Context) error {
	bookId := c.Param("bookId")
	date := c.QueryParam("date")
	dateParsed, err := time.Parse(time.RFC1123, date)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	readingLog, err := s.TrackingService.GetDailyLog(bookId, dateParsed)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// TODO: render view for reading log
	return c.JSON(200, readingLog)
}

func (s *TrackingHandler) GetByUserBookId(c echo.Context) error {
	id := c.Param("id")
	progress, err := s.TrackingService.GetByUserBookId(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(200, progress)
}

func (s *TrackingHandler) Create(c echo.Context) error {
	progress := &models.ReadingProgress{}
	if err := c.Bind(progress); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	err := s.TrackingService.Create(progress)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, progress)
}

func (s *TrackingHandler) Update(c echo.Context) error {
	progress := &models.ReadingProgress{}
	if err := c.Bind(progress); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	err := s.TrackingService.Update(progress)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, progress)
}

func (s *TrackingHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	err := s.TrackingService.Delete(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusNoContent)
}
