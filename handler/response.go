package handler

import (
	"github.com/bSkracic/delta-rest/model"
	"github.com/bSkracic/delta-rest/utils"
)

type userResponse struct {
	User struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Token    string `json:"token"`
	} `json:"user"`
}

func newUserResponse(u *model.User) *userResponse {
	r := new(userResponse)
	r.User.Username = u.Username
	r.User.Email = u.Email
	r.User.Token = utils.GenerateJWT(u.ID)
	return r
}

type submissionResposne struct {
	Submission struct {
		ID           uint   `json:"id"`
		MainFileText string `json:"main_file_text"`
		Language     struct {
			ID   uint
			Name string `json:"name"`
		} `json:"language"`
		Author struct {
			Username string `json:"username"`
		} `json:"author"`
	} `json:"submission"`
}

func newSubmissionResponse(s *model.Submission) *submissionResposne {
	r := new(submissionResposne)
	r.Submission.ID = s.ID
	r.Submission.Language.ID = s.Language.ID
	r.Submission.Language.Name = s.Language.Name
	r.Submission.Author.Username = s.User.Username
	r.Submission.MainFileText = string(s.MainFile)
	return r
}

type execEntryResponse struct {
	ID         uint   `json:"id"`
	Status     string `json:"status"`
	Submission struct {
		ID           uint   `json:"id"`
		MainFileText string `json:"main_file_text"`
		Language     struct {
			ID   uint
			Name string `json:"name"`
		} `json:"language"`
		Author struct {
			Username string `json:"username"`
		} `json:"author"`
	} `json:"submission"`
	ExecConfig struct {
		CompilerOpt string `json:"compiler_opt"`
		TimeLimit   uint   `json:"time_limit"`
		MemoryLimit uint   `json:"memory_limit"`
	} `json:"exec_config"`
}

func newExecEntryResponse(e *model.ExecEntry) *execEntryResponse {
	r := new(execEntryResponse)
	r.ID = e.ID
	r.Status = e.Status.String()
	r.Submission.ID = e.Submission.ID
	r.Submission.MainFileText = string(e.Submission.MainFile)
	r.Submission.Language.ID = e.Submission.Language.ID
	r.Submission.Language.Name = e.Submission.Language.Name
	r.Submission.Author.Username = e.Submission.User.Username
	r.ExecConfig.CompilerOpt = e.CompilerOpt
	r.ExecConfig.TimeLimit = e.TimeLimit
	r.ExecConfig.MemoryLimit = e.MemoryLimit
	return r
}
