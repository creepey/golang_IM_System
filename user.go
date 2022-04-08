package main

import (
	"fmt"
	"net"
)

type User struct {
	Name          string
	Conn          net.Conn
	C             chan string
	broadcastchan chan string
}

//监听C管道,将C管道中内容发送给对应客户端
func (t *User) listenC() {
	for {
		m := <-t.C
		t.Conn.Write([]byte(m + "\n"))
	}
}

//接收客户端消息
func (t *User) Reader() {
	buf := make([]byte, 4096)
	for {
		n, err := t.Conn.Read(buf)
		if err != nil || n == 0 {
			fmt.Printf(t.Name + "退出\n")
			t.broadcastchan <- string(t.Name + "退出\n")
			return
		}
		t.broadcastchan <- t.Name + ": " + string(buf[:n-1])
	}
}

func NewUser(conn net.Conn, broadcastchan chan string) (user *User) {
	user = &User{
		Name:          conn.RemoteAddr().String(),
		Conn:          conn,
		C:             make(chan string),
		broadcastchan: broadcastchan,
	}
	broadcastchan <- user.Name + "已连接"
	return
}
