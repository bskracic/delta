package handler

import (
	"github.com/bSkracic/delta-rest/model"
	"github.com/labstack/echo/v4"
)

type userRegisterRequest struct {
	User struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	} `json:"user"`
}

func (r *userRegisterRequest) bind(c echo.Context, u *model.User) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	u.Username = r.User.Username
	u.Email = r.User.Email
	h, err := u.HashPassword(r.User.Password)
	if err != nil {
		return err
	}
	u.Password = h
	return nil
}

type userLoginRequest struct {
	User struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	} `json:"user"`
}

func (r *userLoginRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	return nil
}

type newExecEntryRequest struct {
	SubmissionID uint `json:"submission_id"`
	ExecConfig   struct {
		CompilerOpt string `json:"compiler_opt"`
		TimeLimit   uint   `json:"time_limit"`
		MemoryLimit uint   `json:"memory_limit"`
	} `json:"exec_config"`
}

func (r *newExecEntryRequest) bind(c echo.Context, e *model.ExecEntry) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	// if err := c.Validate(r); err != nil {
	// 	return err
	// }

	e.SubmissionId = r.SubmissionID
	e.CompilerOpt = r.ExecConfig.CompilerOpt
	e.TimeLimit = r.ExecConfig.TimeLimit
	e.MemoryLimit = r.ExecConfig.MemoryLimit
	return nil
}

type newSubmissionAndExecEntry struct {
	MainFileText string `json:"main_file_text"`
	LanguageID   uint   `json:"language_id" validate:"numeric"`
	ExecConfig   struct {
		CompilerOpt string `json:"compiler_opt"`
		TimeLimit   uint   `json:"time_limit"`
		MemoryLimit uint   `json:"memory_limit"`
	} `json:"exec_config"`
}

func (r *newSubmissionAndExecEntry) bind(c echo.Context, s *model.Submission, e *model.ExecEntry) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	s.LanguageID = r.LanguageID
	s.MainFile = []byte(r.MainFileText)
	e.CompilerOpt = r.ExecConfig.CompilerOpt
	e.TimeLimit = r.ExecConfig.TimeLimit
	e.MemoryLimit = r.ExecConfig.MemoryLimit
	return nil
}
