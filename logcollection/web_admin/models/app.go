package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/jmoiron/sqlx"
)

type AppInfo struct {
	AppId       int    `db:"app_id" json:"appId""`
	AppName     string `db:"app_name" json:"appName"`
	AppType     string `db:"app_type" json:"appType"`
	CreateTime  string `db:"create_time" json:"createTime"`
	DevelopPath string `db:"dev_path" json:"developPath"`
	IP          []string
}

var (
	Db *sqlx.DB
)

func InitDb(db *sqlx.DB) {
	Db = db
}

func GetAllAppInfo() (appList []AppInfo, err error) {
	err = Db.Select(&appList, "select app_id, app_name, app_type, create_time, dev_path from tbl_app_info")
	if err != nil {
		logs.Warn("Get All App Info failed, err:%v", err)
		return
	}
	return
}

func CreateApp(info *AppInfo)(err error){

	tx, err := Db.Begin()
	if err != nil{
		logs.Warn("createApp failed, Db.Begin error: %v",err)
		return
	}

	//事务提交或者回滚
	defer func(){
		if err != nil{
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	result, err := tx.Exec("insert into tbl_app_info(app_name,app_type,dev_path) values (?,?,?)",
		info.AppName, info.AppType, info.DevelopPath)
	if err != nil{
		logs.Warn("insert into tbl_app_info failed ,err:%v",err)
		return
	}

	id, err := result.LastInsertId()
	if err != nil{
		logs.Warn("get lastInsertId failed, err:%v",err)
		return
	}
	for _,ip := range info.IP{
		_, err = tx.Exec("insert into tbl_app_ip (app_id,ip) values (?,?)", id, ip)
		if err != nil {
			logs.Warn("insert into app_tbl_ip failed, err%v",err)
		}
		return
	}

	return
}


func GetIPInfoByName(appName string) (iplist []string, err error) {

	var appId []int
	err = Db.Select(&appId, "select app_id from tbl_app_info where app_name=?", appName)
	if err != nil || len(appId) == 0 {
		logs.Warn("select app_id failed, Db.Exec error:%v", err)
		return
	}

	err = Db.Select(&iplist, "select ip from tbl_app_ip where app_id=?", appId[0])
	if err != nil {
		logs.Warn("Get All App Info failed, err:%v", err)
		return
	}

	logs.Warn("get appId ",iplist)
	return
}