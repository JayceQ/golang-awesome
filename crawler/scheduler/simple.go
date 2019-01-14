package scheduler

import "golang-awesome/crawler/engine"

//简单的并发处理，所有工作协程，公用一个输入通道
type SimpleScheduler struct {
	WorkerChan chan engine.Request
}

func (s *SimpleScheduler) Run() {
	s.WorkerChan = make(chan engine.Request)
}

func (s *SimpleScheduler) Submit(request engine.Request) {
	go func() {
		s.WorkerChan <- request
	}()
}

func (s *SimpleScheduler) GetWorkChan() chan engine.Request {
	return s.WorkerChan
}

func (s *SimpleScheduler) WorkReady(chan engine.Request) {
	panic("implement me")
}
