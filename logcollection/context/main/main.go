package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Result struct {
	r   *http.Response
	err error
}

func process() {
	timeout, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	c := make(chan Result, 1)
	request, e := http.NewRequest("GET", "http://google.com", nil)
	if e != nil {
		fmt.Println("http request failed, err:", e)
		return
	}
	go func() {
		response, i := client.Do(request)
		pack := Result{r: response, err: i}
		c <- pack
	}()

	select {
	case <-timeout.Done():
		tr.CancelRequest(request)
		<-c
		fmt.Println("time out!")
	case res := <-c:
		defer res.r.Body.Close()
		bytes, _ := ioutil.ReadAll(res.r.Body)
		fmt.Printf("Server response : %s ", bytes)
	}
	return
}

func main() {
	process()
}
