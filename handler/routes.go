package handler

import (
	"github.com/bSkracic/delta-rest/router/middleware"
	"github.com/bSkracic/delta-rest/utils"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(v1 *echo.Group) {

	jwtMiddleware := middleware.JWT(utils.JWTSecret)

	// DELETE THIS
	test := v1.Group("/test")
	test.GET("/submissions/:id", h.GetSubmission)
	test.POST("/submissions", h.CreateSubmission)
	// DELETE THIS

	// Submissions
	submissions := v1.Group("/submissions", jwtMiddleware)
	submissions.POST("", h.CreateSubmission)
	submissions.GET("/:id", h.GetSubmission)
	// Get all submissions
	// Change submissions
	// Delete submission

	// ExecEntries
	execEntries := v1.Group("/exec_entries")
	execEntries.POST("/start", h.StartExecution)
	execEntries.GET("/:id", h.GetExecutionEntry)

	// POST start execution(id submission)
	// POST submit && execute (given text or form file) -> how does this affect multiple file problem

}
