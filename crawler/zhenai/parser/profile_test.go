package parser

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"testing"
)

func TestReg(t *testing.T){
	body, err := ioutil.ReadFile("profile_test_data.html")
	if err != nil {
		log.Printf("read file err, %v",err)
		return
	}

	//utf8Reader := transform.NewReader(bytes.NewReader(body),simplifiedchinese.GBK.NewDecoder())
	//body1, _ := ioutil.ReadAll(utf8Reader)

	profileReg := regexp.MustCompile(`<div class="m-btn purple" data-v-bff6f798="">([^<]+)</div>`)
	submatch := profileReg.FindAllSubmatch(body,-1)

	var IdRe = regexp.MustCompile(`<div class="id" data-v-5b109fc3="">ID：([^<]+)</div>`)
	findSubmatch := IdRe.FindSubmatch(body)
	fmt.Printf("%v \n",string(findSubmatch[1]))
	for _, m := range submatch{
		fmt.Printf("%v \n",string(m[1]))
	}
}
