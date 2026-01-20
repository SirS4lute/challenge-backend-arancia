package httpapi

import (
	"errors"
	"net/http"

	"challenge-backend-arancia/internal/application/todos"
	"challenge-backend-arancia/internal/domain"
	"challenge-backend-arancia/internal/ports"

	"github.com/gin-gonic/gin"
)

type todoHandler struct {
	svc *todos.Service
}

type todoResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type createTodoRequest struct {
	Title string `json:"title" binding:"required"`
}

type updateTodoRequest struct {
	Title     string `json:"title" binding:"required"`
	Completed *bool  `json:"completed" binding:"required"`
}

func (h todoHandler) register(r gin.IRoutes) {
	r.GET("/todos", h.list)
	r.POST("/todos", h.create)
	r.PUT("/todos/:id", h.update)
	r.DELETE("/todos/:id", h.delete)
}

func (h todoHandler) list(c *gin.Context) {
	out, err := h.svc.List(c.Request.Context())
	if err != nil {
		writeError(c, err)
		return
	}
	resp := make([]todoResponse, 0, len(out))
	for _, td := range out {
		resp = append(resp, toResponse(td))
	}
	c.JSON(http.StatusOK, resp)
}

func (h todoHandler) create(c *gin.Context) {
	var req createTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	td, err := h.svc.Create(c.Request.Context(), req.Title)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, toResponse(td))
}

func (h todoHandler) update(c *gin.Context) {
	id := c.Param("id")

	var req updateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Completed == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	td, err := h.svc.Update(c.Request.Context(), id, req.Title, *req.Completed)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, toResponse(td))
}

func (h todoHandler) delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func toResponse(td domain.Todo) todoResponse {
	return todoResponse{ID: td.ID, Title: td.Title, Completed: td.Completed}
}

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidTitle):
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid title"})
	case errors.Is(err, ports.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	case errors.Is(err, ports.ErrConflict):
		c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}
}
