package main

import (
	"github.com/julienschmidt/httprouter"
	"golang-awesome/video_server/api/session"
	"log"
	"net/http"
)

type middleWareHandler struct {
	r *httprouter.Router
}

func newMiddleWareHandler(r *httprouter.Router) *middleWareHandler {
	return &middleWareHandler{r: r}
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r * http.Request){
	validateUserSession(r)
	m.r.ServeHTTP(w,r)
}


func RegisterHandlers() *httprouter.Router{
	log.Printf("preparing to post requesr\n")
	r := httprouter.New()
	r.POST("/user",CreateUser)
	r.POST("/user/:username",Login)
	r.GET("/user/:username",GetUserInfo)
	r.POST("/user/:username/videos",AddNewVideo)
	r.GET("/user/:username/videos",ListAllVideos)
	r.DELETE("/user/:username/videos/:vid-id",DeleteVideo)
	r.POST("/videos/:vid-id/comments",PostComment)
	r.GET("/videos/:vid-id/comments",ShowComments)
	return r
}

func prepare(){
	session.LoadSessionsFromDB()
}

func main(){
	prepare()
	r := RegisterHandlers()
	m := newMiddleWareHandler(r)
	http.ListenAndServe(":8000",m)
}