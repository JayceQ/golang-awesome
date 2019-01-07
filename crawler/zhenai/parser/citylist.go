package parser

import (
	"fmt"
	"golang-awesome/crawler/engine"
	"regexp"
)

const cityListReg = `<a href="(http://www.zhenai.com.zhenghun/\w+)"[^>]*>([^<]+)</a>`

func ParseCityList(constents []byte) engine.ParserResult{
	re := regexp.MustCompile(cityListReg)
	matches := re.FindAllSubmatch(constents,-1)
	results := engine.ParserResult{}

	for i, m :=range matches{
		results.Items = append(results.Items, "City "+string(m[2]))
		results.Requests  = append(results.Requests, engine.Request{
			Url:string(m[1]),
			ParserFunc:ParseCity,
		})

		if i == 10{
			fmt.Printf("only testing %d pages, then break", i)
		}

	}
	return results
}