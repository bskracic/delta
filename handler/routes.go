package handler

import "github.com/labstack/echo/v4"

func (h *Handler) Register(v1 *echo.Group) {
	submissions := v1.Group("/submissions")

	submissions.POST("", h.CreateSubmission)
	submissions.GET("/:id", h.GetSubmission)
}
