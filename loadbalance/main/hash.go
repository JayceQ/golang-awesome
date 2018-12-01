package main

import (
	"golang-awesome/loadbalance/balance"
	"fmt"
	"math/rand"
	"hash/crc32"
)

func init(){
	balance.RegisterBalance("hash",&Hashbalance{})
}

type Hashbalance struct {
}

func (p *Hashbalance) DoBalance(insts []*balance.Instance, key ...string) (inst *balance.Instance, err error) {

	var defKey  = fmt.Sprintf("%d", rand.Int())
	if len(key) > 0 {
		defKey = key[0]
	}

	lens := len(insts)
	if lens == 0 {
		err = fmt.Errorf("no backend balance")
		return
	}

	crcTable := crc32.MakeTable(crc32.IEEE)
	hashVal := crc32.Checksum([]byte(defKey), crcTable)
	index := int(hashVal) % lens
	inst = insts[index]

	return
}
