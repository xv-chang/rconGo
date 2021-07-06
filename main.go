package main

import (
	"fmt"

	"github.com/xv-chang/rconGo/core"
)

func main() {

	serverHost := "you Server Ip:port"
	rconPassword := "rconpassword"
	//test server query
	sq := core.NewServerQuery(serverHost)
	defer sq.Close()
	r1 := sq.GetPlayers()
	fmt.Println(r1)

	//test rcon
	client := core.NewRCONClient(serverHost, rconPassword)
	defer client.Close()
	err := client.SendAuth()
	if err != nil {
		println("rcon password error")
		return
	}
	r2, err := client.ExecCommand("sm_cvar z_tank_health 8000")
	if err != nil {
		println("no auth")
	}
	fmt.Println(r2)

}
