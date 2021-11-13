package handler

import (
	"github.com/bSkracic/delta-cli/execentry"
	"github.com/bSkracic/delta-cli/submission"
)

type Handler struct {
	submissionStore submission.Store
	execEntryStore  execentry.Store
}

func NewHandler(ss submission.Store, es execentry.Store) *Handler {
	return &Handler{
		submissionStore: ss,
		execEntryStore:  es,
	}
}
