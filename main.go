package main

import (
	"fmt"
	"os"

	"github.com/bSkracic/delta-rest/db"
	"github.com/bSkracic/delta-rest/handler"
	"github.com/bSkracic/delta-rest/lib/dockercli"
	"github.com/bSkracic/delta-rest/model"
	"github.com/bSkracic/delta-rest/router"
	"github.com/bSkracic/delta-rest/store"
)

func main() {

	f, err := os.OpenFile("logfile.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))
	}
	defer f.Close()

	r := router.New(f)
	v1 := r.Group("/api")

	d := db.New()
	d.AutoMigrate(
		&model.User{},
		&model.Language{},
		&model.Submission{},
		&model.ExecEntry{},
	)
	us := store.NewUserStore(d)
	ss := store.NewSubmissionStore(d)
	ls := store.NewLanguageStore(d)
	es := store.NewExecEntryStore(d)
	dc := dockercli.CreateClient()

	h := handler.NewHandler(us, ls, ss, es, *dc)
	h.Register(v1)
	v1.POST("/register", h.SignUp)
	v1.POST("/login", h.Login)
	r.Logger.Fatal(r.Start("127.0.0.1:8080"))
}
