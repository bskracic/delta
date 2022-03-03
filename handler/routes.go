package handler

import (
	"github.com/bSkracic/delta-rest/router/middleware"
	"github.com/bSkracic/delta-rest/utils"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(v1 *echo.Group) {

	public := v1.Group("/public")
	public.POST("/login", h.Login)
	public.POST("/register", h.SignUp)
	public.POST("/exec", h.SubmitAndExecute)
	public.GET("/entries/:id", h.GetPublicExecutionEntry)

	jwtMiddleware := middleware.JWT(utils.JWTSecret)

	submissions := v1.Group("/submissions", jwtMiddleware)
	submissions.POST("", h.CreateSubmission)
	submissions.GET("/:id", h.GetSubmission)
	submissions.GET("/user", h.GetSubmissionForAuthor)
	// Change submissions with :id
	// Delete submission with :id

	execEntries := v1.Group("/exec_entries", jwtMiddleware)
	execEntries.POST("/start", h.StartExecution)
	// execEntries.POST("/file") execute given file; reuse submission request
	execEntries.GET("/:id", h.GetExecutionEntry)
	execEntries.GET("/submission/:id", h.GetExecutionEntriesForSubmission)
}
