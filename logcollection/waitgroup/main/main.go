package main

import (
	"fmt"
	"sync"
	"time"
)

func main(){
	wg := sync.WaitGroup{}

	for i := 0;i<10;i++{
		wg.Add(1)
		go cacl(&wg,i)
	}

	wg.Wait()
	fmt.Println("all groutine finish")
}

func cacl(w *sync.WaitGroup,i int){
	fmt.Println("cacl:",i)
	time.Sleep(time.Second)
	w.Done()
}