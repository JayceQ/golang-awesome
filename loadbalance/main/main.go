package main

import (
	"golang-awesome/loadbalance/balance"
	"fmt"
	"math/rand"
	"time"
	"os"
)

func main() {
	var insts []*balance.Instance

	for i := 0; i < 16; i++ {
		host := fmt.Sprintf("192.168.%d.%d", rand.Intn(255), rand.Intn(255))
		one := balance.NewInstance(host, 8080)
		insts = append(insts, one)
	}

	var balanceName = "random"

	if len(os.Args) > 1 {
		balanceName = os.Args[1]
	}

	for {
		inst, err := balance.DoBalance(balanceName, insts)
		if err != nil {
			fmt.Println("do banlance err:", err)
		}

		fmt.Println(inst)
		time.Sleep(time.Second)
	}
}
