package engine

// 并发的引擎
// 引擎将请求发送给调度器，调度器纷发给workers, workers的结果再返回给引擎
// 所有worker 引用一个源
type ConcurrentEngine struct {
	MaxWorkerCount int
	Scheduler      Scheduler
	ItemChan       chan Item
	RequestWorker Processor
}
type Processor func(request Request) (ParseResult, error)
type Scheduler interface {
	Submit(request Request)
	GetWorkerChan() chan Request

	Run()
	Ready
}
type Ready interface {
	WorkerReady(chan Request)
}

func (e *ConcurrentEngine) Run(seed ...Request) {
	out := make(chan ParseResult, 1024)
	e.Scheduler.Run()

	for i := 0; i < e.MaxWorkerCount; i++ {
		e.createWorker(e.Scheduler.GetWorkerChan(), out, e.Scheduler)
	}
	for _, r := range seed {
		if IsDuplicate(r.Url) {
			continue
		}
		e.Scheduler.Submit(r)
	}
	//itemCount := 0
	for {
		result := <-out
		for _, item := range result.Items {
			//itemCount++
			//log.Warn("Got Item: #%d %v", itemCount, item)
			go func() { e.ItemChan <- item }()
		}
		for _, r := range result.Requests {
			if IsDuplicate(r.Url) {
				continue
			}
			e.Scheduler.Submit(r)
		}
	}

}
func (e *ConcurrentEngine) createWorker(in chan Request, out chan ParseResult, s Ready) {
	go func() {
		for {
			s.WorkerReady(in)
			request := <-in
			result, err := e.RequestWorker(request)
			if err != nil {
				continue
			}
			out <- result
		}
	}()
}