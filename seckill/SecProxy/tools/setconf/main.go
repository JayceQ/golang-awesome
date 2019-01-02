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

type SecInfoConf struct {
	ProductId int
	StartTime int
	EndTime   int
	Status    int
	Total     int
	Left      int
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

	var SecInfoConfArr []SecInfoConf
	SecInfoConfArr = append(
		SecInfoConfArr,
		SecInfoConf{
			ProductId: 1029,
			StartTime: 1505008800,
			EndTime:   1505012400,
			Status:    0,
			Total:     1000,
			Left:      1000,
		},
	)
	SecInfoConfArr = append(
		SecInfoConfArr,
		SecInfoConf{
			ProductId: 1027,
			StartTime: 1505008800,
			EndTime:   1505012400,
			Status:    0,
			Total:     2000,
			Left:      1000,
		},
	)

	data, err := json.Marshal(SecInfoConfArr)
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
