package handler

import (
	"github.com/bSkracic/delta-rest/interface/execentry"
	language "github.com/bSkracic/delta-rest/interface/languageconf"
	"github.com/bSkracic/delta-rest/interface/submission"
	"github.com/bSkracic/delta-rest/interface/user"
	"github.com/bSkracic/delta-rest/lib/dockercli"
)

type Handler struct {
	userStore       user.Store
	languageStore   language.Store
	submissionStore submission.Store
	execEntryStore  execentry.Store
	docker          dockercli.Dockercli
}

func NewHandler(us user.Store, ls language.Store, ss submission.Store, es execentry.Store, d dockercli.Dockercli) *Handler {
	return &Handler{
		userStore:       us,
		languageStore:   ls,
		submissionStore: ss,
		execEntryStore:  es,
		docker:          d,
	}
}
