package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type HomePage struct {
	Name string
}

type UserPage struct {
	Name string
}

func homeHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cname, err := r.Cookie("username")
	sid, err := r.Cookie("session")
	if err != nil {
		p := &HomePage{Name: "JAYCE"}
		t, e := template.ParseFiles("./templates/home.html")
		if e != nil {
			log.Printf("Parsing template home.html error: %s", e)
			return
		}
		t.Execute(w, p)
		return
	}

	if len(cname.Value) != 0 && len(sid.Value) != 0 {
		http.Redirect(w, r, "/userhome", http.StatusFound)
		return
	}
}

func userHomeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cname, err := r.Cookie("username")
	_, err = r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fname := r.FormValue("username")
	var p *UserPage
	if len(cname.Value) != 0 {
		p = &UserPage{Name: cname.Value}
	} else if len(fname) != 0 {
		p = &UserPage{Name: fname}
	}
	files, err := template.ParseFiles("./templates/userhome.html")
	if err != nil {
		log.Printf("parsing userhome.html error: %s", err)
		return
	}

	files.Execute(w, p)
}

func apiHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method != http.MethodPost {
		re, _ := json.Marshal(ErrorRequestNotRecognized)
		io.WriteString(w, string(re))
		return
	}

	res, _ := ioutil.ReadAll(r.Body)
	apiBody := &ApiBody{}
	if err := json.Unmarshal(res, apiBody); err != nil {
		re, _ := json.Marshal(ErrorRequestBodyParseFailed)
		io.WriteString(w, string(re))
		return
	}
	request(apiBody, w, r)
	defer r.Body.Close()
}

func proxyHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	parse, _ := url.Parse("http://127.0.0.1:9000/")
	proxy := httputil.NewSingleHostReverseProxy(parse)
	proxy.ServeHTTP(w, r)
}
