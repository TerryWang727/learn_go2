package main

import (
	"bufio"
	"log"
	"net"
	"fmt"
)

type Message struct {
	MsgChan chan string
}

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:1235")
	if err != nil {
		log.Fatalf("listen error:%v\n", err)
	}
	fmt.Println("程序启动,开始监听1235端口...")
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("accept error:%v\n", err)
			continue
		}
		//启动读conn的协程
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	rd := bufio.NewReader(conn)

	msg := &Message{make(chan string, 8)}
	go sendMsg(conn, msg.MsgChan)

	for {
		line, _, err := rd.ReadLine()
		if err != nil {
			log.Printf("read error:%v\n", err)
			return
		}

		msg.MsgChan <- string(line)
	}
}

func sendMsg(conn net.Conn, ch <-chan string) {
	wr := bufio.NewWriter(conn)

	for msg := range ch {
		wr.WriteString("hello")
		wr.WriteString(msg)
		wr.Flush()
	}
}
