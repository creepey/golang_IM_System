package main

import (
	"fmt"
	"net"
)

type User struct {
	Name   string
	Conn   net.Conn
	C      chan string
	Server *Server
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
			t.Exit()
			return
		}
		t.Server.Broadcast(t, string(buf[:n-1]))
		// else if string(buf[:n-1]) == "list" {
		// 	for key := range t.Server.OlineMap {
		// 		t.Conn.Write([]byte(key + "\n"))
		// 	}

		// } else {

		// }

	}
}

//用户上线
func (t *User) Login() {
	//将用加入全局用户表中
	t.Server.OlineMap[t.Conn.RemoteAddr().String()] = t
	t.Server.Broadcast(t, t.Name+"上线")
}

//用户下线
func (t *User) Exit() {
	delete(t.Server.OlineMap, t.Conn.RemoteAddr().String())
	fmt.Println(t.Name + "连接断开")
	t.Server.Broadcast(t, t.Name+"下线")
}

func NewUser(conn net.Conn, server *Server) (user *User) {
	user = &User{
		Name:   conn.RemoteAddr().String(),
		Conn:   conn,
		C:      make(chan string),
		Server: server,
	}
	user.Login()
	return
}
