package routers

import (
	"github.com/astaxie/beego"
	"golang-awesome/logcollection/web_admin/controllers/AppController"
)

func init() {
	beego.Router("/index", &AppController.AppController{}, "*:AppList")
	beego.Router("/app/list", &AppController.AppController{}, "*:AppList")
}
