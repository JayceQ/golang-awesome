package service

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func convertLogLevel(level string) int  {
	switch level{
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace
	default:
		return logs.LevelDebug
	}
}

func InitLogger()(err error){

	config := make(map[string]interface{})
	config["filename"] = AppConfig.LogPath
	config["level"] = convertLogLevel(AppConfig.LogLevel)

	bytes, err := json.Marshal(config)
	if err != nil{
		fmt.Println("initLogger failed, marshal err: ",err)
		return
	}

	logs.SetLogger(logs.AdapterFile,string(bytes))
	logs.SetLogger(logs.AdapterConsole,string(bytes))
	beego.SetLogFuncCall(true)
	return
}