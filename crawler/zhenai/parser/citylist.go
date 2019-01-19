package parser

import (
	"golang-awesome/crawler/engine"
	"regexp"
)

const cityListReg = `<a href="(http://www.zhenai.com.zhenghun/\w+)"[^>]*>([^<]+)</a>`

func ParseCityList(contents []byte,_ string) engine.ParserResult{
	results := engine.ParserResult{}

	reg := regexp.MustCompile(cityListReg)
	matches := reg.FindAllSubmatch(contents,-1)


	for _, m :=range matches{
		results.Requests = append(results.Requests, engine.Request{
			Url: string(m[1]),
			Parse:engine.NewFuncParser(ParseCity, "ParseCity"),
		})
	}
	return results
}