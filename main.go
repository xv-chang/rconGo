package main

import (
	"fmt"

	"github.com/xv-chang/rconGo/core"
)

func main() {

	// sq := core.NewServerQuery("myl4d2.tk:27015")
	sq := core.NewServerQuery("43.241.50.78:32063")
	defer sq.Close()
	r := sq.GetPlayers()
	fmt.Println(r)

}
