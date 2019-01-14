package engine

import (
	"github.com/astaxie/beego/logs"
	"golang-awesome/crawler/fetcher"
)

func Worker(r Request) (ParserResult,error){
	body, err := fetcher.Fetch(r.Url)
	if err != nil {
		logs.Error("请求[%s]失败，%s",r.Url,err)
		return ParserResult{},err
	}
	return r.Parse.Parse(body,r.Url),nil
}
