package service

import (
	"fmt"
	"github.com/astaxie/beego/logs"
)

var (
	secKillConf *SecKillConf
)
func InitService(serviceConf *SecKillConf){
	secKillConf = serviceConf
	logs.Debug("load service success, config:%v",serviceConf)
}

func SecInfo(productId int) (data map[string]interface{},code int,err error){
	secKillConf.RWSecProductLock.RLock()
	defer secKillConf.RWSecProductLock.RUnlock()

	 v,ok := secKillConf.SecProductInfoMap[productId]
	 if !ok {
	 	code = ErrNotFoundProductId
	 	err = fmt.Errorf("not fount product id: %d",productId)
	 	return
	 }

	 data = make(map[string]interface{})
	 data["productId"] = productId
	 data["startTime"] = v.StartTime
	 data["endTime"] = v.EndTime
	 data["status"] = v.Status
	 return
}
