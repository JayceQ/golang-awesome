package main

import (
	"github.com/julienschmidt/httprouter"
	"golang-awesome/video_server/scheduler/taskrunner"
	"net/http"
)

func RegisterHandlers() *httprouter.Router {
	r := httprouter.New()
	r.GET("/video-delete-record/:vid-id", vidDelRecHandler)
	return r
}

func main() {
	go taskrunner.Start()
	r := RegisterHandlers()
	http.ListenAndServe(":9001",r)
}
