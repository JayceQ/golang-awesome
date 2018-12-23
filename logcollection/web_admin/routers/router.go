package routers

import (
	"github.com/astaxie/beego"
	"golang-awesome/logcollection/web_admin/controllers/AppController"
	"golang-awesome/logcollection/web_admin/controllers/LogController"
)

func init() {
	beego.Router("/index", &AppController.AppController{}, "*:AppList")

	beego.Router("/app/list", &AppController.AppController{}, "*:AppList")
	beego.Router("/app/apply", &AppController.AppController{}, "*:AppApply")
	beego.Router("/app/create", &AppController.AppController{}, "*:AppCreate")

	beego.Router("/log/list", &LogController.LogController{}, "*:LogList")
	beego.Router("/log/apply", &LogController.LogController{}, "*:LogApply")
	beego.Router("/log/create", &LogController.LogController{}, "*:LogCreate")


	beego.Router("/hello", &AppController.AppController{}, "*:TestRestful")
}
