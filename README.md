# rconGo

golang version RCON Protocal

[Server Queries](https://developer.valvesoftware.com/wiki/Server_queries)

[Source RCON Protocal](https://developer.valvesoftware.com/wiki/Source_RCON_Protocol)

## How to Use

> Server Queries

```
    serverHost := "you Server Ip:port"
   	sq := core.NewServerQuery(serverHost)
	defer sq.Close()
	info := sq.GetInfo()
	fmt.Println(info)
	players := sq.GetPlayers()
	fmt.Println(players)
	rules := sq.GetRules()
	fmt.Println(rules)
```

> RCON

```
    serverHost := "you Server Ip:port"
	rconPassword := "rconpassword"
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

```
