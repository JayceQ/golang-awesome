package session

import (
	"fmt"
	"golang-awesome/video_server/api/dbops"
	"golang-awesome/video_server/api/defs"
	"golang-awesome/video_server/api/utils"
	"sync"
	"time"
)

var sessionMap *sync.Map

func init(){
	sessionMap = &sync.Map{}
}

//返回毫秒数
func nowInMilli() int64{
	return time.Now().UnixNano()/1000000
}

func deleteExpireSession(sid string){
	sessionMap.Delete(sid)
	dbops.DeleteSession(sid)
}

func LoadSessionsFromDB() *sync.Map{
	r, err := dbops.RetrieveAllSessions()
	if err != nil {
		return nil
	}
	r.Range(func(key, value interface{}) bool {
		ss := value.(*defs.SimpleSession)
		sessionMap.Store(key,ss)
		return true
	})
	return sessionMap
}

func GenerateNewSessionId(uname string) string{
	uuid, _ := utils.NewUUID()
	ct := time.Now().UnixNano()/1000000
	ttl := ct + 30 * 60 * 1000
	ss := &defs.SimpleSession{UserName:uname, TTL:ttl}
	sessionMap.Store(uuid, ss)
	err := dbops.InsertSession(uuid, ttl, uname)
	if err != nil {
		return fmt.Sprintf("Error of GenerateNewSessionId: %s",err)
	}
	return uuid
}

func IsSessionExpired(sid string)(string,bool){
	ss, ok := sessionMap.Load(sid)
	if ok{
		ct := nowInMilli()
		if ss.(*defs.SimpleSession).TTL < ct {
			deleteExpireSession(sid)
			return "",true
		}
		return ss.(*defs.SimpleSession).UserName, false
	}
	return "",true
}