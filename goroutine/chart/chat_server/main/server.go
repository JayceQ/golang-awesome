package main

import (
	"fmt"
	"net"
)

func runServer(addr string) (err error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("listen failed, ", err)
		return
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("accept failed, ", err)
			continue
		}

		go proccess(conn)
	}
}

func proccess(conn net.Conn) {

	defer conn.Close()
	client := &Client{
		conn: conn,
	}

	err := client.Process()
	if err != nil {
		fmt.Println("client process failed, ", err)
		return
	}
}
