package engine

import (
	"github.com/gpmgo/gopm/modules/log"
)
// 单任务版引擎
type SimpleEngine struct {

}
func (e *SimpleEngine) Run(queue ...Request) {
	var count = 0
	for len(queue) > 0 {
		r := queue[0]
		queue = queue[1:]

		results, err := Worker(r)
		if err != nil {
			continue
		}
		for _,r := range results.Requests {
			if IsDuplicate(r.Url) {
				continue
			}
			queue = append(queue, r)
		}
		//queue = append(queue, results.Requests...)
		for _, item := range results.Items {
			count++
			log.Warn("Got Item: $%d %v",count, item)
		}
	}
}

