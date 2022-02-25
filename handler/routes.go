package handler

import (
	"github.com/bSkracic/delta-rest/router/middleware"
	"github.com/bSkracic/delta-rest/utils"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(v1 *echo.Group) {

	jwtMiddleware := middleware.JWT(utils.JWTSecret)

	submissions := v1.Group("/submissions", jwtMiddleware)
	submissions.POST("", h.CreateSubmission)
	submissions.GET("/:id", h.GetSubmission)
	submissions.GET("/user", h.GetSubmissionForAuthor)
	// Get all submissions for user with :id
	// Change submissions with :id
	// Delete submission with :id

	execEntries := v1.Group("/exec_entries", jwtMiddleware)
	execEntries.POST("/start", h.StartExecution)
	execEntries.POST("/text", h.SubmitAndExecute)
	// execEntries.POST("/file") execute given file; reuse submission request
	execEntries.GET("/:id", h.GetExecutionEntry)
	execEntries.GET("/submission/:id", h.GetExecutionEntriesForSubmission)

	// POST submit && execute (given text or form file) -> how does this affect multiple file problem

}
