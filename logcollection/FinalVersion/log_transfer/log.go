package main

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/logs"
)

func convertLogLevel(level string) int {

	switch level {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace
	}

	return logs.LevelDebug
}

func initLogger(logPath string, logLevel string) (err error) {

	config := make(map[string]interface{})
	config["filename"] = logPath
	config["level"] = convertLogLevel(logLevel)

	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("initLogger failed, marshal err:", err)
		return
	}

	logs.SetLogger(logs.AdapterFile, string(configStr))
	return
}
