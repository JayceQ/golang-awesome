package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params){
	io.WriteString(w,"Create User Handler")
}

func UserInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params){
	name := p.ByName("name")
	s := fmt.Sprintf("Login user is %s", name)
	io.WriteString(w,s)
}