package router

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"golang-awesome/SecKill/SecProxy/controller"
)

func init() {
	logs.Debug("enter router main")
	beego.Router("/seckill", &controller.SkillController{}, "*:SecKill")
	beego.Router("/secinfo", &controller.SkillController{}, "*:SecInfo")
}

