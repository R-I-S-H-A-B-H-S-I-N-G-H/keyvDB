package main

import (
	"fmt"
	"net"

	"github.com/R-I-S-H-A-B-H-S-I-N-G-H/keyvDB/internal/handler"
	"github.com/R-I-S-H-A-B-H-S-I-N-G-H/keyvDB/internal/server"
	"github.com/tidwall/redcon"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	rdReader := redcon.NewReader(conn)

	cmdHandler := handler.NewCmdHandler()

	for {
		rmsg, err := rdReader.ReadCommand()
		if err != nil {
			fmt.Println("error redcon:", err.Error())
			return
		}

		resp := cmdHandler.HandleDbCommands(rmsg)
		conn.Write(resp)
	}
}

func main() {
	srv := server.NewServer(":6379")
	err := srv.Start(handleConnection)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
