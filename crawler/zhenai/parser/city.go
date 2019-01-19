package parser

import (
	"golang-awesome/crawler/engine"
	"regexp"
)

const cityReg = `<a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([^<]+)</a`

var (
	profileRe = regexp.MustCompile(`<a href="(http://album.zhenai.com/u/\w+)"[^>]*>([^<]+)</a>`)

	cityUrlRe = regexp.MustCompile(`href="(http://www.zhenai.com/zhenghun/[^"]+)`)
)


func ParseCity(contents []byte,_ string) engine.ParserResult{
	rs := engine.ParserResult{}

	match := profileRe.FindAllSubmatch(contents, -1)
	for _, m := range match {
		rs.Requests = append(rs.Requests, engine.Request{
			Url: string(m[1]),
			Parse:NewProfileParser(string(m[2])),
		})
	}

	//取本页面其他城市链接
	//match = cityUrlRe.FindAllSubmatch(contents, -1)
	//for _, m := range match {
	//	rs.Requests = append(rs.Requests, engine.Request{
	//		Url: string(m[1]),
	//		Parse: engine.NewFuncParser(ParseCity, "ParseCity"),
	//	})
	//}
	return rs
}