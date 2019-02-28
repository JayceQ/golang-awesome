package main

import (
	"github.com/julienschmidt/httprouter"
	"golang-awesome/video_server/scheduler/dbops"
	"net/http"
)

func vidDelRecHandler(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	vid := p.ByName("vid-id")
	if len(vid) == 0 {
		sendResponse(w,400, "video id should not be empty")
	}

	err := dbops.AddVideoDeletionRecord(vid)
	if err != nil {
		sendResponse(w, 500, "Internal server error")
		return
	}
	sendResponse(w,200, "delete video successfully")
	return
}
