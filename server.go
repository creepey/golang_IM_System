package main

import (
	"fmt"
	"net"
)

type Server struct {
	ip             string
	port           int
	Broadcast_chan chan string
	OlineMap       map[string]*User
}

func (t *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", t.ip, t.port))
	if err != nil {
		fmt.Printf("listen err: %v\n", err)
	}
	defer listener.Close()
	//开启广播机制
	go t.Broadcast()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("conn err: %v\n", err)
		}
		fmt.Println(conn.RemoteAddr().String() + "连接已建立")
		go t.Handler(conn)
	}

}

//用户业务逻辑
func (t *Server) Handler(conn net.Conn) {
	//创建用户
	user := NewUser(conn, t.Broadcast_chan)
	//将用加入全局用户表中
	t.OlineMap[user.Conn.RemoteAddr().String()] = user
	//监听user维护的chan
	go user.listenC()
	go user.Reader()
}

func (t *Server) Broadcast() {
	for {
		mes := <-t.Broadcast_chan
		for _, one_user := range t.OlineMap {
			one_user.C <- mes
		}
	}
}

func NewServer(ip string, port int) (server *Server) {
	server = &Server{
		ip:             ip,
		port:           port,
		OlineMap:       make(map[string]*User),
		Broadcast_chan: make(chan string),
	}
	return
}
