package main

import (
	"fmt"
	"net"
)

type Server struct {
	ip       string
	port     int
	OlineMap map[string]*User
}

func (t *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", t.ip, t.port))
	if err != nil {
		fmt.Printf("listen err: %v\n", err)
	}
	defer listener.Close()
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
	user := NewUser(conn, t)
	//监听user维护的chan
	go user.listenC()
	go user.Reader()
}

func command(mes string) bool {
	for _, i := range []string{"list"} {
		if i == mes {
			return true
		}
	}
	return false
}

func (t *Server) Broadcast(user *User, mes string) {
	if command(mes) {
	} else {
		for _, one_user := range t.OlineMap {
			if one_user != user {
				one_user.C <- user.Name + ": " + mes
			}
		}
	}
}

func NewServer(ip string, port int) (server *Server) {
	server = &Server{
		ip:       ip,
		port:     port,
		OlineMap: make(map[string]*User),
	}
	return
}
