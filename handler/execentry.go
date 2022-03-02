package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bSkracic/delta-rest/lib/dockercli"
	"github.com/bSkracic/delta-rest/model"
	"github.com/bSkracic/delta-rest/utils"
	"github.com/labstack/echo/v4"
)

func (h *Handler) StartExecution(c echo.Context) error {
	// bind new start execution request
	var req newExecEntryRequest
	e := new(model.ExecEntry)
	if err := req.bind(c, e); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	s, err := h.submissionStore.GetSubmission(req.SubmissionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	} else if s == nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}

	c.Response().After(func() {
		h.execute(s, e, req.ExecConfig.TimeLimit)
	})

	e.Status = model.Running
	e.Submission = *s
	if err := h.execEntryStore.Create(e); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusCreated, newExecEntryResponse(e))
}

func (h *Handler) GetExecutionEntry(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}

	if e, err := h.execEntryStore.GetExecEntry(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	} else if e == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	} else {
		return c.JSON(http.StatusOK, newExecEntryResponse(e))
	}
}

func (h *Handler) GetExecutionEntriesForSubmission(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}

	var s *model.Submission
	if s, err = h.submissionStore.GetSubmission(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	} else if s == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	if e, err := h.execEntryStore.GetExecEntriesForSubmission(uint(id)); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	} else {
		// var r []*execEntryResponse
		// for _, ent := range e {
		// 	r = append(r, newExecEntryResponse(&ent))
		// }
		return c.JSON(http.StatusOK, newSubmissionAndExecEntriesResponse(s, e))
	}
}

func (h *Handler) SubmitAndExecute(c echo.Context) error {
	var req newSubmissionAndExecEntry
	s := new(model.Submission)
	e := new(model.ExecEntry)

	if err := req.bind(c, s, e); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}

	s.UserID = userIDFromToken(c)

	if err := h.submissionStore.CreateSubmission(s); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	e.Status = model.Running
	e.Submission = *s
	if err := h.execEntryStore.Create(e); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	c.Response().After(func() {
		h.execute(s, e, req.ExecConfig.TimeLimit)
	})

	return c.JSON(http.StatusCreated, newExecEntryResponse(e))
}

func (h *Handler) execute(s *model.Submission, e *model.ExecEntry, timeLimit uint) {

	start := time.Now()
	log.Print("-----------------------------------------")
	// NOTE: Currently restarting containers if they are stopped
	// FIXME: Handle errors on docker operations
	// Should be experimented with start vs restart time performance

	contId := h.docker.RetreiveAvailableContainer(s.Language.Name, s.Language.Image)

	log.Printf("%s Retrieved container: %s\n", fmt.Sprint(e.ID), time.Since(start))
	stage1 := time.Now()

	// Create directory {exec_entry_id} where main file and optional stdin file will be placed
	dir := fmt.Sprint(e.ID)
	if err := h.docker.CreateDir(contId, dir); err != nil {
		log.Print(err)
	}

	log.Printf("%s Created dir: %s\n", fmt.Sprint(e.ID), time.Since(stage1))
	stage2 := time.Now()

	path := fmt.Sprintf("/%s/", dir)
	// path := "/"
	if err := h.docker.Copy(contId, s.Language.MainFileName, s.MainFile, path); err != nil {
		log.Print(err)
	}

	log.Printf("%s Copied main file: %s\n", fmt.Sprint(e.ID), time.Since(stage2))

	// Copy input if exists
	if e.Stdin != "" {
		stage21 := time.Now()
		if err := h.docker.Copy(contId, "stdin.txt", []byte(e.Stdin), path); err != nil {
			log.Print(err)
		}
		log.Printf("%s Copied stdin: %s\n", fmt.Sprint(e.ID), time.Since(stage21))
	}

	ch := make(chan *dockercli.ExecOutput, 1)
	// Compile source code
	if s.Language.CompileCmd != "" {

		stage3 := time.Now()
		cmd := fmt.Sprintf("cd %s && %s 2>&1", dir, s.Language.CompileCmd)

		go h.docker.Exec(contId, cmd, ch)
		eout := <-ch
		if eout.ExitCode != 0 {
			e.Result = eout.Stdout
			e.ExitCode = eout.ExitCode
			e.Status = model.Failed
			h.execEntryStore.Update(e)
			log.Printf("%s Compiled: %s\n", fmt.Sprint(e.ID), time.Since(stage3))
			return
		}
		log.Printf("%s Compiled: %s\n", fmt.Sprint(e.ID), time.Since(stage3))
	}

	stage4 := time.Now()
	// Execute compiled code (or interpret without compilation step)
	cmd := fmt.Sprintf("cd %s && %s 2>&1", dir, s.Language.ExecuteCmd)
	if e.Stdin != "" {
		cmd = fmt.Sprintf("%s < stdin.txt", cmd)
	}

	go h.docker.Exec(contId, cmd, ch)

	// If time limit is set, terminate program after said time
	if timeLimit != 0 {
		select {
		case res := <-ch:
			e.ExitCode = res.ExitCode
			if res.ExitCode != 0 {
				e.Status = model.Failed
				e.Result = res.Stdout
			} else {
				e.Status = model.Finished
				e.Result = res.Stdout
			}
		case <-time.After(time.Duration(timeLimit) * time.Millisecond):
			e.Status = model.Interrupted
			go h.docker.Kill(contId)
		}
	} else {
		res := <-ch
		e.ExitCode = res.ExitCode
		if res.ExitCode != 0 {
			e.Status = model.Failed
			e.Result = res.Stdout
		} else {
			e.Status = model.Finished
			e.Result = res.Stdout
		}
	}

	log.Printf("%s Executed: %s\n", fmt.Sprint(e.ID), time.Since(stage4))

	stage5 := time.Now()
	if err := h.execEntryStore.Update(e); err != nil {
		log.Fatalln(err)
	}

	log.Printf("%s Saved to db: %s\n", fmt.Sprint(e.ID), time.Since(stage5))

	log.Printf("%s Total time: %s\n", fmt.Sprint(e.ID), time.Since(start))
	log.Print("-----------------------------------------")
}
