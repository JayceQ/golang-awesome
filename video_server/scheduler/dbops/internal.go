package dbops

import (
	"log"
)

func ReadVideoDeletionRecord(count int)(ids []string,err error){
	stmtOut, err := dbConn.Prepare("select video_id from video_del_rec limit ?")
	if err != nil {
		return
	}

	rows, err := stmtOut.Query(count)
	if err != nil {
		log.Printf("Query VideoDeletionRecord error: %s",err)
		return
	}

	for{
		var id string
		if err := rows.Scan(&id); err !=nil{
			return
		}
		ids = append(ids, id)
	}

	defer stmtOut.Close()
	return
}

func DelVideoDeletionRecord(vid string) (err error){
	stmtDel, err := dbConn.Prepare("delete from video_del_rec where video_id = ?")

	if err != nil {
		return
	}

	_, err = stmtDel.Exec(vid)
	if err != nil {
		log.Printf("Deleting VideoDeletionRecord error: %s", err)
		return
	}

	defer stmtDel.Close()
	return
}