package dbops

import "log"

func AddVideoDeletionRecord(vid string) error{
	stmtIns, err := dbConn.Prepare("insert into video_del_rec(video_id) values (?)")
	if err != nil{
		return err
	}

	_, err = stmtIns.Exec(vid)
	if err != nil {
		log.Printf("AddVideoDeletionRecord error: %v",err)
		return err
	}

	defer stmtIns.Close()
	return nil
}