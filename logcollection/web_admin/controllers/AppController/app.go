package AppController

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"golang-awesome/logcollection/web_admin/models"
	"strings"
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

func (a *AppController) AppApply(){
	logs.Debug("enter apply controller")
	a.Layout = "layout/layout.html"
	a.TplName = "app/apply.html"
}

func (a *AppController) AppCreate(){
	logs.Debug("enter appCreate controller")
	appName := a.GetString("appName")
	appType := a.GetString("appType")
	developPath := a.GetString("devPath")
	ipstr := a.GetString("iplist")

	a.Layout = "layout/layout.html"

	if len(appName) == 0 || len(appType) == 0 || len(developPath) == 0 || len(ipstr) == 0 {
		a.Data["Error"] = fmt.Sprintf("非法参数")
		a.TplName = "app/error.html"

		logs.Warn("invalid parameter")
		return
	}

	appInfo := &models.AppInfo{}
	appInfo.AppName = appName
	appInfo.AppType = appType
	appInfo.DevelopPath = developPath
	appInfo.IP = strings.Split(ipstr, ",")

	err := models.CreateApp(appInfo)
	if err != nil {
		a.Data["Error"] = fmt.Errorf(err.Error())
		a.TplName = "app/error.html"
		return
	}

	a.Redirect("/app/list", 302)
}


func (a *AppController) TestRestful(){
	appList, err := models.GetAllAppInfo()
	if err != nil{
		logs.Warn("get appInfo failed ,err:",err)
		return
	}
	a.Data["json"] = map[string]interface{}{"success": 200, "data": appList}
	a.ServeJSON()
}
