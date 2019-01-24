package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func RegisterHandlers() *httprouter.Router{
	router := httprouter.New()

	router.POST("/user",CreateUser)
	router.GET("/info/:name",UserInfo)

	return router
}

func main(){
	r := RegisterHandlers()
	http.ListenAndServe(":8000",r)
}