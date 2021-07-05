package core

import (
	"bytes"
	"net"
)

type ServerQuery struct {
	Conn *net.UDPConn
}
type ServerInfo struct {
	Protocol      byte
	Name          string
	Map           string
	Folder        string
	Game          string
	ID            uint16
	Players       byte
	MaxPlayers    byte
	Bots          byte
	ServerType    byte
	Environment   byte
	Visibility    byte
	VAC           byte
	Version       string
	ExtraDataFlag byte
}
type PlayerInfo struct {
	Index    byte
	Name     string
	Score    uint32
	Duration float32
}
type RuleInfo struct {
	Name  string
	Value string
}

func NewServerQuery(host string) *ServerQuery {
	raddr, _ := net.ResolveUDPAddr("udp", host)
	conn, _ := net.DialUDP("udp", nil, raddr)
	return &ServerQuery{Conn: conn}
}

func (query *ServerQuery) GetChallenge() []byte {
	query.Conn.Write([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0x55, 0xFF, 0xFF, 0xFF, 0xFF})
	var buf [9]byte
	query.Conn.Read(buf[0:])
	return buf[5:]
}

func (query *ServerQuery) GetInfo() *ServerInfo {
	challenge := query.GetChallenge()
	b := bytes.Buffer{}
	b.Write([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0x54})
	b.Write([]byte("Source Engine Query"))
	b.WriteByte(0)
	b.Write(challenge)
	query.Conn.Write(b.Bytes())
	readBuffer := make([]byte, 1024)
	_, err := query.Conn.Read(readBuffer)
	if err != nil {
		panic(err)
	}
	offset := 5
	info := &ServerInfo{}
	info.Protocol = ReadUint8(readBuffer, &offset)
	info.Name = ReadString(readBuffer, &offset)
	info.Map = ReadString(readBuffer, &offset)
	info.Folder = ReadString(readBuffer, &offset)
	info.Game = ReadString(readBuffer, &offset)
	info.ID = ReadUInt16(readBuffer, &offset)
	info.Players = ReadUint8(readBuffer, &offset)
	info.MaxPlayers = ReadUint8(readBuffer, &offset)
	info.Bots = ReadUint8(readBuffer, &offset)
	info.ServerType = ReadUint8(readBuffer, &offset)
	info.Environment = ReadUint8(readBuffer, &offset)
	info.Visibility = ReadUint8(readBuffer, &offset)
	info.VAC = ReadUint8(readBuffer, &offset)
	info.Version = ReadString(readBuffer, &offset)
	info.ExtraDataFlag = ReadUint8(readBuffer, &offset)
	return info
}

func (query *ServerQuery) GetPlayers() []PlayerInfo {
	challenge := query.GetChallenge()
	b := bytes.Buffer{}
	b.Write([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0x55})
	b.Write(challenge)
	query.Conn.Write(b.Bytes())
	readBuffer := make([]byte, 4096)
	_, err := query.Conn.Read(readBuffer)
	if err != nil {
		panic(err)
	}
	offset := 5
	count := ReadUint8(readBuffer, &offset)
	players := make([]PlayerInfo, count)
	for i := byte(0); i < count; i++ {
		player := &PlayerInfo{}
		player.Index = ReadUint8(readBuffer, &offset)
		player.Name = ReadString(readBuffer, &offset)
		player.Score = ReadUInt32(readBuffer, &offset)
		player.Duration = ReadFloat32(readBuffer, &offset)
		players[i] = *player
	}
	return players
}
func (query *ServerQuery) GetRules() []RuleInfo {
	challenge := query.GetChallenge()
	b := bytes.Buffer{}
	b.Write([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0x56})
	b.Write(challenge)
	query.Conn.Write(b.Bytes())
	readBuffer := make([]byte, 4096)
	end, err := query.Conn.Read(readBuffer)
	if err != nil {
		panic(err)
	}
	offset := SearchHeader(readBuffer, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x45}, end)
	count := ReadUInt16(readBuffer, &offset)
	rules := make([]RuleInfo, count)
	for i := uint16(0); i < count; i++ {
		rule := &RuleInfo{}
		rule.Name = ReadString(readBuffer, &offset)
		if rule.Name == "" {
			break
		}
		rule.Value = ReadString(readBuffer, &offset)
		rules[i] = *rule
	}
	return rules
}

func (query *ServerQuery) Close() error {
	return query.Conn.Close()
}
