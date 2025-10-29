package main

import user_service "ticketing-infra/rpc-server/user-service"

func main() {
	go user_service.Start()
	select {}
}
