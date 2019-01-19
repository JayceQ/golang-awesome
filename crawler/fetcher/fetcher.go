package fetcher

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// 100毫秒执行一次请求
var rateLimiter = time.Tick(50 * time.Millisecond)

func Fetch(url string) ([]byte, error){
	//<-rateLimiter
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36")
	var httpClient = http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("抓取出错了, 返回码[%d]", resp.StatusCode)
	}
	bufBody := bufio.NewReader(resp.Body)
	utf8Reader := transform.NewReader(bufBody, determineEncoding(bufBody).NewDecoder())
	body, err := ioutil.ReadAll(utf8Reader)
	return body, err
}

func determineEncoding(r *bufio.Reader) encoding.Encoding{
	bytes, err := r.Peek(1024)
	if err != nil{
		log.Printf("Fetcher error: %v",err)
		return unicode.UTF8
	}

	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
