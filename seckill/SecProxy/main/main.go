package main

import (
	"github.com/astaxie/beego"
)

func main(){

	err := InitConfig()
	if err != nil{
		panic(err)
		return
	}

	err = InitSecKill()
	if err != nil{
		panic(err)
		return
	}

	beego.Run()
}