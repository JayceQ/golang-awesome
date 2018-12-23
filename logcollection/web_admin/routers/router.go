package routers

import (
	"github.com/astaxie/beego"
	"golang-awesome/logcollection/web_admin/controllers/AppController"
)

func init() {
	beego.Router("/index", &AppController.AppController{}, "*:AppList")
	beego.Router("/app/list", &AppController.AppController{}, "*:AppList")
	beego.Router("/app/apply", &AppController.AppController{}, "*:AppApply")
	beego.Router("/app/create", &AppController.AppController{}, "*:AppCreate")


	beego.Router("/hello", &AppController.AppController{}, "*:TestRestful")
}
