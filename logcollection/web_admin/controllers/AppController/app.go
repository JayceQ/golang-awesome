package AppController

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"golang-awesome/logcollection/web_admin/models"
)

type AppController struct {
	beego.Controller
}

func (this *AppController) AppList(){
	logs.Debug("enter index controller")

	this.Layout = "layout/layout.html"
	appList, err := models.GetAllAppInfo()
	if err != nil {
		this.Data["Error"] = fmt.Sprintf("服务器繁忙")
		this.TplName = "app/error.html"

		logs.Warn("get app list failed, err:%v", err)
		return
	}

	logs.Debug("get app list succ, data:%v", appList)
	this.Data["applist"] = appList

	this.TplName = "app/index.html"
}
