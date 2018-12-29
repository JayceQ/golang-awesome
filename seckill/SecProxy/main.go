package main

import (
	"golang-awesome/seckill/SecProxy/load"
	_ "golang-awesome/seckill/SecProxy/router"
	"github.com/astaxie/beego"
)

func main(){

	err := load.InitConfig()
	if err != nil{
		panic(err)
		return
	}

	err = load.InitSecKill()
	if err != nil{
		panic(err)
		return
	}

	beego.Run()
}