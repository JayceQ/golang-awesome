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
