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

func (a *AppController) AppList(){
	logs.Debug("enter index controller")

	a.Layout = "layout/layout.html"
	appList, err := models.GetAllAppInfo()
	if err != nil {
		a.Data["Error"] = fmt.Sprintf("服务器繁忙")
		a.TplName = "app/error.html"

		logs.Warn("get app list failed, err:%v", err)
		return
	}

	logs.Debug("get app list succ, data:%v", appList)
	a.Data["applist"] = appList

	a.TplName = "app/index.html"
}

func (a *AppController) TestHello(){
	a.Ctx.WriteString("hello world")
}
