package main

import (
	"log"

	"github.com/shivaraj-shanthaiah/user-management-apigateway/pkg/config"
	"github.com/shivaraj-shanthaiah/user-management-apigateway/pkg/server"
	"github.com/shivaraj-shanthaiah/user-management-apigateway/pkg/user"
)

func main() {
	cnfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Error while loading config file")
	}

	server := server.NewServer()
	user.NewUserRoute(server.R, *cnfg)
	server.StartServer(cnfg.APIPORT)
}
