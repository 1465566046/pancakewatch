package main

import (
	"os"
	"pancakewatch/pcwdb"
	"pancakewatch/route"

	"github.com/gin-gonic/gin"
	"github.com/twilio/twilio-go"
)

func main() {
	r := gin.Default()

	// TODO: make index a template and extract all other
	// frontend files to an assets folder.
	r.StaticFile("/", "./frontend/index.html")
	r.StaticFile("/index.css", "./frontend/index.css")
	r.StaticFile("/logo.png", "./frontend/logo.png")
	r.StaticFile("/favicon.ico", "./frontend/favicon.ico")

	db, err := pcwdb.New("/tmp/pancakewatch.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r.POST("/subscribe", route.Subscribe(db))

	client := twilio.NewRestClientWithParams(twilio.RestClientParams{
		Username:   os.Getenv("USERNAME"),
		Password:   os.Getenv("PASSWORD"),
		AccountSid: os.Getenv("ACCOUNT_SID"),
	})
	go runChecker(client, db)
	r.Run()
}
