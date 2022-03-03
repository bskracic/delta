package handler

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/bSkracic/delta-rest/model"
	"github.com/bSkracic/delta-rest/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (h *Handler) CreateSubmission(c echo.Context) error {

	file, err := c.FormFile("main_file")
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, src); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	lID, err := strconv.Atoi(c.FormValue("language_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}

	uID := userIDFromToken(c)
	b := buf.Bytes()
	s := model.Submission{
		UserID:     &uID,
		LanguageID: uint(lID),
		MainFile:   b,
	}

	if err := h.submissionStore.CreateSubmission(&s); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	return c.JSON(http.StatusCreated, newSubmissionResponse(&s))
}

func (h *Handler) GetSubmission(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}

	if s, err := h.submissionStore.GetSubmission(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	} else if s == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	} else {
		return c.JSON(http.StatusOK, newSubmissionResponse(s))
	}
}

func (h *Handler) GetSubmissionForAuthor(c echo.Context) error {

	uID := userIDFromToken(c)
	if s, err := h.submissionStore.GetSubmissionsForAuthor(uID); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	} else {
		var r []*submissionResposne
		for _, sub := range s {
			r = append(r, newSubmissionResponse(&sub))
		}
		return c.JSON(http.StatusOK, r)
	}

}

func (h *Handler) DeleteSubmission(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}

	if s, err := h.submissionStore.GetSubmission(uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, utils.NotFound())
		}
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	} else {
		if err = h.submissionStore.DeleteSubmission(s); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
		} else {
			return c.JSON(http.StatusOK, nil)
		}
	}
}
