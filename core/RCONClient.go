package core

import (
	"net"
)

type RCONClient struct {
	Conn *net.UDPConn
}

func NewRCONClient(host string) *RCONClient {
	raddr, _ := net.ResolveUDPAddr("tcp", host)
	conn, _ := net.DialUDP("tcp", nil, raddr)
	return &RCONClient{Conn: conn}
}

func (client *RCONClient) Exec(cmd string) {

}

func (client *RCONClient) Close() error {
	return client.Conn.Close()
}
