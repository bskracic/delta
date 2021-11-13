package main

import (
	"github.com/bSkracic/delta-cli/db"
	"github.com/bSkracic/delta-cli/handler"
	"github.com/bSkracic/delta-cli/model"
	"github.com/bSkracic/delta-cli/router"
	"github.com/bSkracic/delta-cli/store"
)

func main() {
	r := router.New()
	v1 := r.Group("/api")

	d := db.New()
	d.AutoMigrate(
		&model.Submission{},
		&model.ExecEntry{},
	)
	ss := store.NewSubmissionStore(d)

	h := handler.NewHandler(ss, nil)
	h.Register(v1)

	r.Logger.Fatal(r.Start("127.0.0.1:7777"))
}
