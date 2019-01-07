package main

import (
	"golang-awesome/crawler/engine"
	"golang-awesome/crawler/zhenai/parser"
)

func main(){
	engine.Run(engine.Request{
		Url: "http://www.zhenai.com/zhenghun",
		ParserFunc: parser.ParseCityList,
	})


}