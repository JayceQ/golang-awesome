package dbops

import (
	"database/sql"
	"golang-awesome/video_server/api/defs"
	"log"
	"strconv"
	"sync"
)

func InsertSession(sid string, ttl int64, uname string)(err error ){
	ttlstr := strconv.FormatInt(ttl, 10)
	stmtIns, err := dbConn.Prepare("insert into sessions (session_id, TTL, login_name) values (?,?,?)")
	if err != nil {
		return
	}
	_, err = stmtIns.Exec(sid, ttlstr, uname)
	defer stmtIns.Close()
	return
}

func RetrieveSession(sid string)(session *defs.SimpleSession, err error){
	stmtOut, err := dbConn.Prepare("select TTL, login_name from sessions where session_id = ?")
	if err != nil {
		return
	}
	var ttl, uname string
	err = stmtOut.QueryRow(sid).Scan(&ttl, &uname)
	if err != nil && err != sql.ErrNoRows{
		return
	}
	if res, err := strconv.ParseInt(ttl,10,64); err == nil {
		session.TTL = res
		session.UserName = uname
	}else {
		return
	}

	defer stmtOut.Close()
	return
}

func RetrieveAllSessions()(m *sync.Map,err error){
	stmtOut, err := dbConn.Prepare("select * from sessions")
	if err != nil {
		log.Printf("%s",err)
		return
	}
	rows, err := stmtOut.Query()
	if err != nil {
		log.Printf("%s",err)
		return
	}

	for rows.Next() {
		var id, ttlstr, login_name string
		if err := rows.Scan(&id, &ttlstr, &login_name); err != nil {
			log.Printf("retrieve sessions error: %s",err)
			break
		}
		if ttl, err := strconv.ParseInt(ttlstr, 10, 64); err != nil {
			ss := &defs.SimpleSession{
				UserName:login_name,
				TTL: ttl,
			}
			m.Store(id, ss)
			log.Printf("session id: %s, ttl: %d",id,ss.TTL)
		}else{
			log.Printf("parse TTL error: %s",err)
			break
		}
	}
	return
}

func DeleteSession(sid string)(err error){
	stmtOut, err := dbConn.Prepare("delete from sessions where session_id =?")
	if err != nil {
		return
	}
	if _, err := stmtOut.Query(sid); err != nil {
		return
	}

	defer stmtOut.Close()
	return
}
