package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func main() {
	start := time.Now().Unix()
	wg := sync.WaitGroup{}
	for i := 1; i < 100; i++ {
		wg.Add(1)
		go httpGet(i, &wg)
	}
	wg.Wait()
	fmt.Printf("the progress exit, use %d s" ,time.Now().Unix() -start)
}
func httpGet(userId int,wg *sync.WaitGroup) {
	url := fmt.Sprintf("http://localhost:9090/seckill?product_id=1002&user_id=%d",userId)
	resp, err :=   http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(bytes))
	wg.Done()
}