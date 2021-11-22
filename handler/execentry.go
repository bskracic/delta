package handler

import (
	"fmt"
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
	var entry model.ExecEntry
	if err := req.bind(c, &entry); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	// save entry ID

	s, err := h.submissionStore.GetSubmission(req.SubmissionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	} else if s == nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}

	c.Response().After(func() {

		fmt.Printf("\nExecution started\n")

		// ### SUBOPTIMAL EXECUTION
		// Should be put into worker pool or give priority to tasks that are not in batch
		// Retreive available suitable container (may be slow because we are iterating over all possible containers with that image)
		// and killing & removing exited ones

		contId := h.docker.RetreiveAvailableContainer(s.Language.Name, s.Language.Image)

		h.docker.Copy(s.MainFile, s.Language.MainFileName, contId)

		ch := make(chan *dockercli.ExecOutput, 1)
		// Compile source code
		if s.Language.CompileCmd != "" {
			cmd := fmt.Sprintf("%s 2>&1", s.Language.CompileCmd)

			go h.docker.Exec(contId, cmd, ch)
			eout := <-ch
			if eout.ExitCode != 0 {
				// update exec entry eout.Stdout
				entry.Result = eout.Stdout
				entry.ExitCode = eout.ExitCode
				entry.Status = model.Failed
				h.execEntryStore.Update(&entry)
				return
			}
		}

		// Execute compiled code
		go h.docker.Exec(contId, s.Language.ExecuteCmd, ch)

		t := req.ExecConfig.TimeLimit
		var res dockercli.ExecOutput
		if t != 0 {
			select {
			case res := <-ch:
				if res.ExitCode != 0 {
					entry.Status = model.Failed
				} else {
					entry.Status = model.Finished
				}
				// fmt.Printf("STDOUT:\n%s\nSTDERR:\n%s\nEXIT CODE: %d\n", res.Stdout, res.Stderr, res.ExitCode)
			case <-time.After(time.Duration(t) * time.Millisecond):
				entry.Status = model.Interrupted
				// fmt.Println("Interrupted :/")
			}
		} else {
			res := <-ch
			if res.ExitCode != 0 {
				entry.Status = model.Failed
			} else {
				entry.Status = model.Finished
			}
			// fmt.Printf("STDOUT:\n%s\nSTDERR:\n%s\nEXIT CODE: %d\n", res.Stdout, res.Stderr, res.ExitCode)
		}

		entry.Result = res.Stdout
		entry.ExitCode = res.ExitCode
		h.execEntryStore.Update(&entry) // Log fatal
	})

	entry.Status = model.Running
	entry.Submission = *s
	if err := h.execEntryStore.Create(&entry); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusCreated, newExecEntryResponse(&entry))
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
