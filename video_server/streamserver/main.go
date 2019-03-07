package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type middleWareHandler struct {
	r *httprouter.Router
	l *ConnLimiter
}

func NewMiddleWareHandler(r *httprouter.Router, cc int) http.Handler {
	return &middleWareHandler{
		r: r,
		l: NewConnLimiter(cc),
	}
}


func RegisterHandlers() *httprouter.Router{
	r := httprouter.New()
	r.GET("/videos/:vid-id",streamHandler)
	r.POST("/upload/:vid-id",uploadHandler)
	r.GET("/test",testPageHandler)
	return r
}

//重写ServeHTTP接口
func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	if !m.l.GetConn() {
		sendErrorResponse(w, http.StatusTooManyRequests, "too many requests")
		return
	}
	m.r.ServeHTTP(w,r)
	defer m.l.ReleaseConn()
}

func main() {
	r := RegisterHandlers()
	m := NewMiddleWareHandler(r, 20)
	http.ListenAndServe(":9000",m)
}