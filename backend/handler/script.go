package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"microclass-backend/model"
	"microclass-backend/service"
	"microclass-backend/store"
)

type Handler struct {
	store     *store.Store
	llmService *service.LLMService
}

func New(s *store.Store, llm *service.LLMService) *Handler {
	return &Handler{store: s, llmService: llm}
}

func (h *Handler) Generate(c *gin.Context) {
	var req model.GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	script, err := h.llmService.GenerateScript(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.store.Save(script); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save script"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"script": script})
}

func (h *Handler) GenerateStream(c *gin.Context) {
	var req model.GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Writer.Flush()

	sendEvent := func(event string, data any) {
		body, _ := json.Marshal(data)
		fmt.Fprintf(c.Writer, "event: %s\ndata: %s\n\n", event, body)
		c.Writer.Flush()
	}

	script, err := h.llmService.GenerateScriptStream(req, func(chunk string) {
		sendEvent("chunk", gin.H{"text": chunk})
	})
	if err != nil {
		sendEvent("error", gin.H{"error": err.Error()})
		return
	}

	if err := h.store.Save(script); err != nil {
		sendEvent("error", gin.H{"error": "failed to save script"})
		return
	}

	sendEvent("done", gin.H{"script": script})
	io.WriteString(c.Writer, ": heartbeat\n\n")
	c.Writer.Flush()
}

func (h *Handler) ListScripts(c *gin.Context) {
	scripts, err := h.store.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"scripts": scripts})
}

func (h *Handler) GetScript(c *gin.Context) {
	id := c.Param("id")
	script, err := h.store.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "script not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"script": script})
}

func (h *Handler) DeleteScript(c *gin.Context) {
	id := c.Param("id")
	if err := h.store.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
