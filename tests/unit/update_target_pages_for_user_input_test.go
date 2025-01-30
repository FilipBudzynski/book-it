package unit

import (
	"errors"
	"testing"
	"time"

	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestUpdateTargetPagesForUserInput(t *testing.T) {
	mockRepo := new(MockProgressRepository)
	progressService := services.NewProgressService(mockRepo)

	progressID := uint(123)
	progressIDString := "123"
	logID := uint(2)

	today := time.Now()
	progress := &models.ReadingProgress{
		ID:         progressID,
		TotalPages: 50,
		EndDate:    today.AddDate(0, 0, 3),
		Completed:  false,
		DailyProgress: []models.DailyProgressLog{
			{ID: 1, Date: today.AddDate(0, 0, -1), PagesRead: 10},
			{ID: 2, Date: today, PagesRead: 0, TargetPages: 10},
			{ID: 3, Date: today.AddDate(0, 0, 1), PagesRead: 0},
		},
	}

	t.Run("Normal case - TargetPages updated", func(t *testing.T) {
		mockRepo.On("Get", progressID).Return(progress, nil)
		mockRepo.On("GetById", progressIDString).Return(progress, nil)
		mockRepo.On("Update", progress).Return(nil)

		updatedProgress, err := progressService.UpdateTargetPagesForUserInput(progressIDString, logID)

		assert.NoError(t, err)
		assert.Equal(t, progress, updatedProgress)

		assert.Equal(t, 10, progress.DailyProgress[0].TargetPages)
		assert.Equal(t, 10, progress.DailyProgress[1].TargetPages)
		assert.Equal(t, 10, progress.DailyProgress[2].TargetPages)
	})

	t.Run("Completed progress - No update", func(t *testing.T) {
		progress.Completed = true

		mockRepo.On("Get", progressID).Return(progress, nil)
		mockRepo.On("GetById", progressIDString).Return(progress, nil)

		updatedProgress, err := progressService.UpdateTargetPagesForUserInput(progressIDString, logID)

		assert.NoError(t, err)
		assert.Equal(t, progress, updatedProgress)

		assert.Equal(t, 10, progress.DailyProgress[0].TargetPages)
		assert.Equal(t, 10, progress.DailyProgress[1].TargetPages)
		assert.Equal(t, 10, progress.DailyProgress[2].TargetPages)
	})
}

func TestUpdateTargetPagesForUserInput_RepoError(t *testing.T) {
	mockRepo := new(MockProgressRepository)
	progressService := services.NewProgressService(mockRepo)
	progressIDString := "123"
	logID := uint(2)

	t.Run("Repo error - Get returns nil", func(t *testing.T) {
		mockRepo.On("GetById", progressIDString).Return(&models.ReadingProgress{}, errors.New("not found"))
		updatedProgress, err := progressService.UpdateTargetPagesForUserInput(progressIDString, logID)

		assert.Error(t, err)
		assert.Nil(t, updatedProgress)
		assert.Equal(t, "not found", err.Error())
	})
}
