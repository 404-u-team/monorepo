package handlers

import (
	"net/http"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type testDataHandler struct {
	svc services.TestDataService
}

func NewTestDataHandler(svc services.TestDataService) *testDataHandler {
	return &testDataHandler{svc: svc}
}

// Start launches asynchronous test-data generation.
// GET /generate-test-data
func (h *testDataHandler) Start(c *gin.Context) {
	if err := h.svc.Start(); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "генерация тестовых данных запущена"})
}

// Status returns the current generation progress.
// GET /generate-test-data/status
func (h *testDataHandler) Status(c *gin.Context) {
	c.JSON(http.StatusOK, h.svc.Status())
}

// Cancel stops the running generation job.
// GET /generate-test-data/cancel
func (h *testDataHandler) Cancel(c *gin.Context) {
	h.svc.Cancel()
	c.JSON(http.StatusOK, gin.H{"message": "задание на генерацию данных отменено"})
}
