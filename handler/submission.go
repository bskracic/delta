package handler

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/bSkracic/delta-cli/errutil"
	"github.com/bSkracic/delta-cli/model"
	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateSubmission(c echo.Context) error {

	file, err := c.FormFile("mainFile")
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, struct{ err string }{err: "Ubij se"})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, struct{ err string }{err: "Ubij se"})
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, src); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, struct{ err string }{err: "Ubij se"})
	}

	b := buf.Bytes()
	s := model.Submission{Language: c.FormValue("language"), MainFile: b}

	if err := h.submissionStore.CreateSubmission(&s); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, struct{ err string }{err: "Ubij se"})
	}

	return c.JSON(http.StatusCreated, struct{ id uint }{id: s.ID})
}

func (h *Handler) GetSubmission(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, struct {
			message string
		}{message: "id should be provided"})
	}

	if s, err := h.submissionStore.GetSubmission(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	} else if s == nil {
		return c.JSON(http.StatusNotFound, errutil.NotFound())
	} else {
		var buf bytes.Buffer
		buf.Write(s.MainFile)
		return c.JSON(http.StatusOK, struct {
			Submission   model.Submission `json:"submission"`
			MainFileText string           `json:"main_file_text"`
		}{
			Submission:   *s,
			MainFileText: buf.String(),
		})
	}
}

// submit and execute
