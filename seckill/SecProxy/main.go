package main

import "github.com/astaxie/beego"

func main(){

	err := initConfig()
	if err != nil{
		panic(err)
		return
	}

	err = initSecKill()
	if err != nil{
		panic(err)
		return
	}

	beego.Run()
}