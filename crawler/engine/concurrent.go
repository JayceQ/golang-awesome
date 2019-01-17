package engine

//并发的引擎
//引擎将请求发送给调度器，调度器分发给workers，workers的结果再返回给引擎
//所有的worker统一调度
type Scheduler interface{
	Run()
	Submit(request Request)
	GetWorkChan() chan Request
	Ready
}

type Ready interface {
	WorkReady(chan Request)
}

type ConcurrentEngine struct {
	MaxWorkerCount int
	Scheduler Scheduler
	ItemChan chan Item
	RequestWorker Processor
}

type Processor func(request Request)(ParserResult,error)

func (c *ConcurrentEngine) Run(seed ...Request){
	out := make(chan ParserResult, 1024)
	c.Scheduler.Run()

	for i :=0; i < c.MaxWorkerCount; i++{
		c.createWorker(c.Scheduler.GetWorkChan(),out,c.Scheduler)
	}

	for _, r := range seed{
		if IsDuplicate(r.Url){
			continue
		}
		c.Scheduler.Submit(r)
	}

	for {
		result := <-out
		for _, item:=  range result.Items{
			go func() {c.ItemChan <- item}()
		}

		for _, r := range result.Requests{
			if IsDuplicate(r.Url){
				continue
			}
			c.Scheduler.Submit(r)
		}
	}
}

func (c *ConcurrentEngine) createWorker(in chan Request,out chan ParserResult,r Ready){
	go func() {
		for {
			r.WorkReady(in)
			request := <- in
			result, err := c.RequestWorker(request)
			if err != nil {
				continue
			}
			out <- result
		}
	}()
}