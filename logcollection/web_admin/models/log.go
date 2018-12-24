package models

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	etcdclient "go.etcd.io/etcd/clientv3"
	"golang-awesome/logcollection/logagent/tail"
	"time"
)

var(
	etcdClient *etcdclient.Client
)
type LogInfo struct {
	AppId      int    `db:"app_id" json:"app_id"`
	AppName    string `db:"app_name"`
	LogId      int    `db:"log_id"`
	Status     int    `db:"status"`
	CreateTime string `db:"create_time"`
	LogPath    string `db:"log_path"`
	Topic      string `db:"topic"`
}

func InitEtcd(client *etcdclient.Client) {
	etcdClient = client
}


func GetAllLogInfo() (loglist []LogInfo, err error) {
	err = Db.Select(&loglist,
		"select a.app_id, b.app_name, a.create_time, a.topic, a.log_id, a.status, a.log_path from tbl_log_info a, tbl_app_info b where a.app_id=b.app_id")
	if err != nil {
		logs.Warn("Get All App Info failed, err:%v", err)
		return
	}
	return
}

func CreateLog(info *LogInfo) (err error){
	tx, err := Db.Begin()

	if err != nil{
		logs.Warn("create log failed,err %v",err)
		return
	}

	defer func(){
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	var appList []AppInfo
	err = Db.Select(&appList, "select app_id,app_name from tbl_app_info where app_name = ?", info.AppName)

	if err != nil || len(appList) == 0{
		logs.Warn("select appInfo failed, err:%v",err)
		return
	}
	info.AppId = appList[0].AppId
	info.AppName = appList[0].AppName
	logs.Debug(info.AppId)
	result, err := tx.Exec("insert into tbl_log_info (app_id,app_name,log_path,topic) values (?,?,?,?)",
		info.AppId, info.AppName,info.LogPath, info.Topic)

	if err != nil {
		logs.Warn("insert into tbl_log_info falied, err %v",err)
		return
	}

	_, err = result.LastInsertId()

	if err != nil{
		logs.Warn("get last insert id failes, err %v",err)
		return
	}

	return
}

func SetLogConfToEtcd(etcdKey string,info *LogInfo)(err error){

	var logConfArr []tail.CollectConf

	logConfArr = append(
		logConfArr,
		tail.CollectConf{
			LogPath:info.LogPath,
			Topic:info.Topic,
		},
	)
	logs.Warn("etcdctl:",etcdKey)
	bytes, err := json.Marshal(logConfArr)
	if err != nil {
		logs.Warn("marshal log conf failed, err%v",err)
		return
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	_, err = etcdClient.Put(ctx, etcdKey, string(bytes))
	cancelFunc()

	if err != nil {
		logs.Warn("put log conf to etcd failed, err%v",err)
		return
	}

	logs.Debug("put etcd success, data:%v",string(bytes))

	return
}