package core

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"time"
)

type RCONClient struct {
	Conn     *net.TCPConn
	Password string
	Authed   bool
	LastID   int32
}
type RCONPacket struct {
	Size int32
	ID   int32
	Type int32
	Body string
}

const (
	RESPONSE_VALUE = 0
	AUTH_RESPONSE  = 2
	EXECCOMMAND    = 2
	AUTH           = 3
)

func NewRCONClient(host string, password string) *RCONClient {
	raddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		panic(err)
	}
	conn, err := net.DialTCP("tcp", nil, raddr)
	//set timeout 2s
	conn.SetDeadline(time.Now().Add(time.Duration(time.Second * 2)))
	if err != nil {
		panic(err)
	}

	client := &RCONClient{Conn: conn, Password: password, LastID: 0}

	return client
}

func (client *RCONClient) ExecCommand(command string) (string, error) {
	if !client.Authed {
		return "", errors.New("no auth")
	}

	client.LastID++
	client.SendPacket(&RCONPacket{
		ID:   client.LastID,
		Type: EXECCOMMAND,
		Body: command,
	})
	client.SendEmptyPacket()
	packets := client.ReadData()
	r := ""
	for _, v := range packets {
		r += v.Body
	}
	return r, nil
}

func (client *RCONClient) Close() error {
	return client.Conn.Close()
}

func (client *RCONClient) ReadPacket() *RCONPacket {
	readBuffer := make([]byte, 4096)
	_, err := client.Conn.Read(readBuffer)
	if err != nil {
		println(err)
	}
	offset := 0
	return &RCONPacket{
		Size: ReadInt32(readBuffer, &offset),
		ID:   ReadInt32(readBuffer, &offset),
		Type: ReadInt32(readBuffer, &offset),
		Body: ReadString(readBuffer, &offset),
	}
}

func (client *RCONClient) SendPacket(packet *RCONPacket) {
	body := []byte(packet.Body)
	b := bytes.Buffer{}
	packet.Size = int32(len(body) + 9)
	binary.Write(&b, binary.LittleEndian, packet.Size)
	binary.Write(&b, binary.LittleEndian, packet.ID)
	binary.Write(&b, binary.LittleEndian, packet.Type)
	b.Write(body)
	b.WriteByte(0)
	client.Conn.Write(b.Bytes())
}

func (client *RCONClient) SendAuth() error {
	client.LastID++
	client.SendPacket(&RCONPacket{
		ID:   client.LastID,
		Type: AUTH,
		Body: client.Password,
	})
	client.ReadAuthData()
	if !client.Authed {
		return errors.New("auth failed")
	}
	return nil
}

func (client *RCONClient) SendEmptyPacket() {
	client.SendPacket(&RCONPacket{
		ID:   client.LastID,
		Type: RESPONSE_VALUE,
		Body: "",
	})
}
func (client *RCONClient) ReadAuthData() {
	_ = client.ReadPacket()
	p := client.ReadPacket()
	if p.Type == AUTH_RESPONSE {
		if p.ID > -1 {
			client.Authed = true
		}
	}
}

func (client *RCONClient) ReadData() []*RCONPacket {
	packets := make([]*RCONPacket, 1)
	p := client.ReadPacket()
	packets[0] = p
	if p.Type == RESPONSE_VALUE {
		for p.Body != "" {
			p = client.ReadPacket()
			packets = append(packets, p)
		}
	}
	return packets
}
