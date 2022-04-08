package main

import "net"

type User struct {
	Name string
	Conn net.Conn
	C    chan string
}

//监听C管道,将C管道中内容发送给对应客户端
func (t *User) listenC() {
	for {
		m := <-t.C
		t.Conn.Write([]byte(m + "\n"))
	}
}

func NewUser(conn net.Conn, broadcastchan chan string) (user *User) {
	user = &User{
		Name: conn.RemoteAddr().String(),
		Conn: conn,
		C:    make(chan string),
	}
	broadcastchan <- user.Name + "已连接"
	return
}
