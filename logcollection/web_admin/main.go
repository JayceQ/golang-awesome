package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/jmoiron/sqlx"
	"golang-awesome/logcollection/web_admin/models"
	_ "golang-awesome/logcollection/web_admin/routers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego"
)

func initDb() (err error) {
	database, err := sqlx.Open("mysql", "root:P@ssw0rd@tcp(182.61.137.53:3306)/logagent")
	if err != nil {
		logs.Warn("open mysql failed,", err)
		return
	}

	models.InitDb(database)
	return
}
func main() {
	err := initDb()

	if err != nil {
		logs.Warn("initDb failed, err:%v", err)
		return
	}
	beego.Run()
}

