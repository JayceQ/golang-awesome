package main

import (
	"fmt"
	"github.com/astaxie/beego/config"
)

func main() {
	conf, e := config.NewConfig("ini", "./logagent.conf")
	if e != nil {
		fmt.Println("new config failed ,err :", e)
		return
	}

	port, e := conf.Int("server::port")
	if e != nil {
		fmt.Println("read server:port failed,err:",e)
		return
	}

	fmt.Println("port:",port)

}
