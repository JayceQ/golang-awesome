package main

import (
	"context"
	"encoding/json"
	"fmt"
	etcdClient "go.etcd.io/etcd/clientv3"
	"time"


)

const (
	EtcdKey = "/net/badme/backend/seckill/product"
)

type SecProductInfoConf = struct {
	ProductId int
	StartTime int64
	EndTime int64
	Status int
	Total int
	Left int
	OnePersonBuyLimit int
	BuyRate float64
	//每秒组多能卖多少个
	SoldMaxLimit int
}


func SetLogConfToEtcd() {
	cli, err := etcdClient.New(etcdClient.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("connect failed, err:", err)
		return
	}

	fmt.Println("connect succ")
	defer cli.Close()

	now := time.Now().Unix()
	var SecProductInfoConfArr []SecProductInfoConf
	SecProductInfoConfArr = append(
		SecProductInfoConfArr,
		SecProductInfoConf{
			ProductId: 1001,
			StartTime: now - 60,
			EndTime:   now + 3600,
			Status:    0,
			Total:     10,
			OnePersonBuyLimit:1,
			SoldMaxLimit:10000,
			BuyRate: 0.5,
		},
	)
	SecProductInfoConfArr = append(
		SecProductInfoConfArr,
		SecProductInfoConf{
			ProductId: 1002,
			StartTime: now - 60,
			EndTime:   now + 3600,
			Status:    0,
			Total:     10,
			OnePersonBuyLimit:1,
			SoldMaxLimit:10000,
			BuyRate: 0.5,
		},
	)

	data, err := json.Marshal(SecProductInfoConfArr)
	if err != nil {
		fmt.Println("json failed, ", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//cli.Delete(ctx, EtcdKey)
	//return
	_, err = cli.Put(ctx, EtcdKey, string(data))
	cancel()
	if err != nil {
		fmt.Println("put failed, err:", err)
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, EtcdKey)
	cancel()
	if err != nil {
		fmt.Println("get failed, err:", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
}

func main() {
	SetLogConfToEtcd()
}
