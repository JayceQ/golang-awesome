package LogController

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"golang-awesome/logcollection/web_admin/models"
)

type LogController struct {
	beego.Controller
}

func (a *LogController) LogList() {
	logs.Debug("enter logList controller")
	a.Layout = "layout/layout.html"
	a.Layout = "layout/layout.html"
	logList, err := models.GetAllLogInfo()
	if err != nil {
		a.Data["Error"] = fmt.Sprintf("服务器繁忙")
		a.TplName = "app/error.html"

		logs.Warn("get app list failed, err:%v", err)
		return
	}

	logs.Debug("get log list succ, data:%v", logList)
	a.Data["loglist"] = logList

	a.TplName = "log/index.html"
}

func (a *LogController) LogApply() {
	logs.Debug(" enter logApply controller")
	a.Layout = "layout/layout.html"
	a.TplName = "log/apply.html"

}

func (a *LogController) LogCreate() {
	logs.Debug("enter logCreate controller")
	appName := a.GetString("appName")
	logPath := a.GetString("logPath")
	topic := a.GetString("topic")

	a.Layout = "layout/layout.html"
	if len(appName) == 0 || len(logPath) == 0 || len(topic) == 0 {
		a.Data["Error"] = fmt.Sprintf("非法参数")
		a.TplName = "log/error.html"

		logs.Warn("invalid parameter")
		return
	}

	logInfo := &models.LogInfo{
		AppName: appName,
		LogPath: logPath,
		Topic:   topic,
	}

	err := models.CreateLog(logInfo)
	logs.Debug(err)
	if err != nil {
		a.Data["Error"] = fmt.Errorf(err.Error())
		logs.Warn("create logInfo failed, err%v", err)
		a.TplName = "app/error.html"
		return
	}

	iplist, err := models.GetIPInfoByName(appName)
	if err != nil || len(iplist) == 0{
		a.Data["Error"] = fmt.Sprintf("获取项目ip失败，数据库繁忙")
		a.TplName = "log/error.html"

		logs.Warn("invalid parameter")
		return
	}
	keyFormat := "/net/badme/logagent/config/%s"

	logs.Debug(" get iplist:",iplist)

	for _, ip := range iplist {
		key := fmt.Sprintf(keyFormat, ip)
		err = models.SetLogConfToEtcd(key, logInfo)
		if err != nil {
			logs.Warn("Set log conf to etcd failed, err:%v", err)
			continue
		}
	}
	a.Redirect("/log/list", 302)

}
