package scheduler

import "golang-awesome/crawler/engine"

type Scheduler interface{
	Run()
	Submit(request engine.Request)
	GetWorkChan() chan engine.Request
	Ready
}

type Ready interface {
	WorkReady(chan engine.Request)
}