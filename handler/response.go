package handler

import (
	"time"

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
			ID   uint   `json:"id"`
			Name string `json:"name"`
		} `json:"language"`
		Author struct {
			Username string `json:"username"`
		} `json:"author"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"submission"`
}

func newSubmissionResponse(s *model.Submission) *submissionResposne {
	r := new(submissionResposne)
	r.Submission.ID = s.ID
	r.Submission.CreatedAt = s.CreatedAt
	r.Submission.Language.ID = s.Language.ID
	r.Submission.Language.Name = s.Language.Name
	r.Submission.Author.Username = s.User.Username
	r.Submission.MainFileText = string(s.MainFile)
	return r
}

type execEntryResponse struct {
	ExecEntry struct {
		ID         uint      `json:"id"`
		Status     string    `json:"status"`
		Result     string    `json:"result"`
		ExitCode   int       `json:"exit_code"`
		CreatedAt  time.Time `json:"created_at"`
		Submission struct {
			ID           uint      `json:"id"`
			MainFileText string    `json:"main_file_text"`
			CreatedAt    time.Time `json:"created_at"`
			Language     struct {
				ID   uint   `json:"id"`
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
			Stdin       string `json:"stdin"`
		} `json:"exec_config"`
	} `json:"exec_entry"`
}

func newExecEntryResponse(e *model.ExecEntry) *execEntryResponse {
	r := new(execEntryResponse)
	r.ExecEntry.ID = e.ID
	r.ExecEntry.Status = e.Status.String()
	r.ExecEntry.Result = e.Result
	r.ExecEntry.ExitCode = e.ExitCode
	r.ExecEntry.Submission.ID = e.Submission.ID
	r.ExecEntry.CreatedAt = e.CreatedAt
	r.ExecEntry.ExecConfig.Stdin = e.Stdin
	r.ExecEntry.Submission.CreatedAt = e.Submission.CreatedAt
	r.ExecEntry.Submission.MainFileText = string(e.Submission.MainFile)
	r.ExecEntry.Submission.Language.ID = e.Submission.Language.ID
	r.ExecEntry.Submission.Language.Name = e.Submission.Language.Name
	r.ExecEntry.Submission.Author.Username = e.Submission.User.Username
	r.ExecEntry.ExecConfig.CompilerOpt = e.CompilerOpt
	r.ExecEntry.ExecConfig.TimeLimit = e.TimeLimit
	r.ExecEntry.ExecConfig.MemoryLimit = e.MemoryLimit
	return r
}

type submissionAndExecEntriesResponse struct {
	Submission struct {
		ID           uint   `json:"id"`
		MainFileText string `json:"main_file_text"`
		Language     struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		} `json:"language"`
		Author struct {
			Username string `json:"username"`
		} `json:"author"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"submission"`
	ExecEntries []*execEntryResponse `json:"exec_entries"`
}

func newSubmissionAndExecEntriesResponse(s *model.Submission, e []model.ExecEntry) *submissionAndExecEntriesResponse {
	r := new(submissionAndExecEntriesResponse)
	r.Submission.ID = s.ID
	r.Submission.MainFileText = string(s.MainFile)
	r.Submission.Language.ID = s.Language.ID
	r.Submission.Language.Name = s.Language.Name
	r.Submission.Author.Username = s.User.Username
	r.Submission.CreatedAt = s.CreatedAt
	var re []*execEntryResponse
	for _, ent := range e {
		re = append(re, newExecEntryResponse(&ent))
	}
	r.ExecEntries = re
	return r
}
