package service

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"time"
)

var (
	secKillConf *SecKillConf
)

func NewSecRequest() (secRequest *SecRequest) {
	secRequest = &SecRequest{
		ResultChan: make(chan *SecResult, 1),
	}
	return
}

func SecInfoList() (data []map[string]interface{}, code int, err error) {
	secKillConf.RWBlackLock.RLock()
	defer secKillConf.RWBlackLock.RUnlock()

	for _, v := range secKillConf.SecProductInfoMap {
		item, _, err := SecInfoById(v.ProductId)
		if err != nil {
			logs.Error("get product_id[%d] failed, err:%v", err)
			continue
		}

		logs.Debug("get product[%d], result[%v], all[%v] v[%v]", v.ProductId, item, secKillConf.SecProductInfoMap, v)
		data = append(data, item)
	}

	return
}

func SecInfoById(productId int) (data map[string]interface{}, code int, err error) {
	secKillConf.RWBlackLock.RLock()
	defer secKillConf.RWBlackLock.RUnlock()

	v, ok := secKillConf.SecProductInfoMap[productId]
	if !ok {
		code = ErrNotFoundProductId
		err = fmt.Errorf("not found product id: %d", productId)
		return
	}

	start := false
	end := false
	status := "success"

	now := time.Now().Unix()
	if now-v.StartTime < 0 {
		start = false
		end = false
		status = "sec kill is not start"
		code = ErrActiveNotStart
	}

	if now-v.StartTime > 0 {
		start = true
	}

	if now-v.EndTime > 0 {
		start = false
		end = true
		status = "sec kill is already end"
		code = ErrActiveAlreadyEnd
	}

	if v.Status == ProductStatusForceSaleOut || v.Status == ProductStatusSaleOut {
		start = false
		end = true
		status = "product is sale out"
		code = ErrActiveSaleOut
	}

	data = make(map[string]interface{})
	data["product_id"] = productId
	data["start"] = start
	data["end"] = end
	data["status"] = status

	return
}

func SecInfo(productId int) (data []map[string]interface{}, code int, err error) {
	secKillConf.RWSecProductLock.RLock()
	defer secKillConf.RWSecProductLock.RUnlock()

	item, code, err := SecInfoById(productId)

	if err != nil {
		return
	}
	data = append(data, item)
	return
}

func SecKill(req *SecRequest) (data map[string]interface{}, code int, err error) {

	//secKillConf.RWSecProductLock.RLock()
	//defer secKillConf.RWSecProductLock.RUnlock()

	//err = antiSpam(req)
	//if err != nil {
	//	code = ErrUserServiceBusy
	//	logs.Warn("userId[%d] invalid, check failed, req[%v]", req.UserId, req)
	//	return
	//}

	data, code, err = SecInfoById(req.ProductId)
	if code != 0 {
		logs.Warn("userId[%d] secInfoBy id failed, req[%v]", req.UserId, req)
		return
	}

	userKey := fmt.Sprintf("%s_%s", req.UserId, req.ProductId)
	secKillConf.RWSecProductLock.RLock()
	secKillConf.UserConnMap[userKey] = req.ResultChan
	secKillConf.RWSecProductLock.RUnlock()

	secKillConf.SecReqChan <- req
	ticker := time.NewTicker(time.Second * 10)

	defer func() {
		ticker.Stop()
		secKillConf.UserConnMapLock.Lock()
		delete(secKillConf.UserConnMap, userKey)
		secKillConf.UserConnMapLock.Unlock()
	}()

	select {
	case <-ticker.C:
		code = ErrProcessTimeout
		err = fmt.Errorf("request timeout")
		return

	case <-req.CloseNotify:
		code = ErrClientClosed
		err = fmt.Errorf("client already closed")
		return
	case result := <-req.ResultChan:
		code = result.Code
		data["product_id"] = result.ProductId
		data["token"] = result.Token
		data["user_id"] = result.UserId
		return
	}
	return
}
