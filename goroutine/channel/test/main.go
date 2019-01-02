package main

import (
	"fmt"
	"time"
)

func main(){
	msg := make(chan string, 1)

	for i := 0; i < 16; i++ {
		go send(msg,i)
	}

	for i := 0; i < 16; i++ {
		go recv(msg, i)
	}

	time.Sleep(time.Hour)
}

func send(msg chan string,i int) {
	str := fmt.Sprintf("this msg is from %dth goroutine, ",i)
	for {
		msg <- str
		time.Sleep(time.Second)
	}
}

func recv(msgChan chan string, i int) {

	//for msg := range msgChan {
	//	fmt.Printf("%sand receivced by %dth goroutine\n", msg,i)
	//}

	//for{
	//	msg, ok := <- msgChan
	//	if !ok {
	//		break
	//	}
	//	fmt.Printf("%sand receivced by %dth goroutine\n", msg,i)
	//}

	for{
		select {
		case msg := <-msgChan:
			fmt.Printf("%sand receivced by %dth goroutine\n", msg,i)
		}
	}

}
