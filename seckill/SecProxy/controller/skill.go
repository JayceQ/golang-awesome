package controller

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"golang-awesome/seckill/SecProxy/service"
)

type SkillController struct {
	beego.Controller
}

func (p *SkillController) SecKill(){
	p.Data["json"] = "seckill"
	p.ServeJSON()
}

func (p *SkillController) SecInfo(){

	productId ,err := p.GetInt("product_id")
	result := make(map[string]interface{})
	logs.Debug("product_id:%d",productId)

	result["code"] = 0
	result["message"] = "success"

	defer func(){
		p.Data["json"] = result
		p.ServeJSON()
	}()

	if err != nil {
		result["code"] = 1001
		result["message"] = "invalid produce_id"
		logs.Error("invalid request ,get product_id falied, err:%v",err)
		return
	}

	data, code, err := service.SecInfo(productId)
	if err != nil {
		result["code"] = code
		result["message"] = err.Error()
		logs.Error("invalid request,get product_id falied, err :%v",err)
		return
	}

	result["data"] = data
}