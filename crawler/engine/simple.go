package engine

import "github.com/astaxie/beego/logs"

type SimpleEngine struct {

}

func (e *SimpleEngine) Run(queue ...Request){
	var count = 0
	for len (queue) >0 {
		r := queue[0]
		queue = append(queue[1:])

		results, err := Worker(r)
		if err != nil{
			continue
		}
		for _, r := range results.Requests{
			if IsDuplicate(r.Url){
				continue
			}
			queue = append(queue,r)
		}

		for _, item :=range results.Items{
			count ++
			logs.Warn("Got Item: $%d %v",count,item)
		}
	}
}
