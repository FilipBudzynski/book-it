package handlers

import (
	"fmt"
	"net/http"
	"strconv"

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
	UpdateTargetPages(progressID string, logID uint) error
	Delete(id string) error

	// log methods
	GetLog(id string) (*models.DailyProgressLog, error)
	UpdateLog(id string, pagesRead int, comment string) (*models.DailyProgressLog, error)
	RefreshTargetPagesForNewDay(progressID string, logID uint) error
}

type progressHandler struct {
	progressService ProgressService
	userBookService UserBookService
}

func NewProgressHandler(s ProgressService, u UserBookService) *progressHandler {
	return &progressHandler{
		progressService: s,
		userBookService: u,
	}
}

func (h *progressHandler) RegisterRoutes(app *echo.Echo) {
	group := app.Group("/progress")
	group.Use(utils.CheckLoggedInMiddleware) // middleware for protected routes
	group.POST("", h.Create)
	group.GET("/:id", h.GetByUserBookId)
	group.DELETE("/:id", h.Delete)
	// progress log endpoints
	group.Use(utils.CheckLoggedInMiddleware)
	group.PUT("/log/:id", h.UpdateLog)
	// htmx routes
	group.GET("/log/details/modal/:id", h.GetLogModal)
	group.GET("/details/:id", h.GetProgressDetails)
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

	// progressID := fmt.Sprintf("%d", progress.ID)
	// if err := h.progressService.RefreshTargetPagesForNewDay(progressID); err != nil {
	// 	return errs.HttpErrorInternalServerError(err)
	// }

	return utils.RenderView(c, webProgress.ProgressStatistics(progress))
}

func (h *progressHandler) GetProgressDetails(c echo.Context) error {
	id := c.Param("id")
	progress, err := h.progressService.GetByUserBookId(id)
	if err != nil {
		c.Response().Status = http.StatusNotFound
		return errs.HttpErrorNotFound(err)
	}
	userBook, err := h.userBookService.Get(id)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}
	return utils.RenderView(c, webProgress.CardProgress(progress, userBook))
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

func (h *progressHandler) UpdateLog(c echo.Context) error {
	id := c.Param("id")
	comment := c.FormValue("comment")
	pagesRead, err := strconv.Atoi(c.FormValue("pages-read"))
	if err != nil {
		return errs.HttpErrorBadRequest(err)
	}

	log, err := h.progressService.UpdateLog(id, pagesRead, comment)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	progressID := fmt.Sprintf("%d", log.ReadingProgressID)
	if err := h.progressService.UpdateTargetPages(progressID, log.ID); err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	progress, err := h.progressService.GetProgressAssosiatedWithLogId(id)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	if !progress.IsFinishedOnLastLog(log.Date) {
		_ = toast.Info(models.ErrProgressLastDayNotFinished.Error()).SetHXTriggerHeader(c)
	}

	if progress.Completed {
		_ = toast.Success(c, CompletedBookMessage)
	}

	userBook, err := h.userBookService.Get(fmt.Sprintf("%d", progress.UserBookID))
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}

	return utils.RenderView(c, webProgress.ProgressDetailsOverview(progress, userBook))
}

func (h *progressHandler) GetLogModal(c echo.Context) error {
	id := c.Param("id")
	log, err := h.progressService.GetLog(id)
	if err != nil {
		return errs.HttpErrorInternalServerError(err)
	}
	return utils.RenderView(c, webProgress.ProgressLogModal(*log))
}
