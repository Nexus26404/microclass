package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"microclass-backend/config"
	"microclass-backend/model"
)

func (h *Handler) GetSettings(c *gin.Context) {
	mockMode, _, openAIURL, openAIModel, _ := h.store.GetSettings()
	c.JSON(http.StatusOK, gin.H{
		"settings": model.Settings{
			MockMode:    mockMode,
			OpenAIURL:   openAIURL,
			OpenAIModel: openAIModel,
		},
	})
}

func (h *Handler) UpdateSettings(c *gin.Context) {
	var req model.SettingsUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	cfg := config.Get()
	cfg.Update(req.MockMode, req.OpenAIKey, req.OpenAIURL, req.OpenAIModel)

	h.store.SaveSettings(cfg.MockMode, cfg.OpenAIKey, cfg.OpenAIURL, cfg.OpenAIModel)

	c.JSON(http.StatusOK, gin.H{
		"settings": model.Settings{
			MockMode:    cfg.MockMode,
			OpenAIURL:   cfg.OpenAIURL,
			OpenAIModel: cfg.OpenAIModel,
		},
	})
}
