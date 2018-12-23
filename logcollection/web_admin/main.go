package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jmoiron/sqlx"
	"go.etcd.io/etcd/clientv3"
	"golang-awesome/logcollection/web_admin/models"
	_ "golang-awesome/logcollection/web_admin/routers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego"
	"time"
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

func initEtcd()(err error){
	client, e := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})

	if e != nil {
		fmt.Println("connect to etcd failed,err:",e)
		return
	}
	models.InitEtcd(client)
	return
}

func main() {
	err := initDb()

	if err != nil {
		logs.Warn("initDb failed, err:%v", err)
		return
	}

	err = initEtcd()
	if err != nil {
		logs.Warn("initEtcd failed, err%v",err)
		return
	}
	beego.Run()
}

