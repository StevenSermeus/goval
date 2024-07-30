package networking

import (
	"net"

	"github.com/StevenSermeus/goval/src/cache"
	"github.com/StevenSermeus/goval/src/commands"

	"github.com/StevenSermeus/goval/src/config"
	"github.com/StevenSermeus/goval/src/types"
)

func HandleClient(conn net.Conn, cache *cache.Cache, serverConfig *config.Config) {
	client := types.Client{
		Conn:         conn,
		Cache:        cache,
		ServerConfig: serverConfig,
	}
	for {
		command, err := client.Receive()
		if err != nil {
			break
		}
		//This allow non blocking execution of commands on the server
		//So errors might occur since client waiting for response from request 1 might receive response from request 2 first and vice versa
		go commands.ExecuteCommad(&client, command)
	}

	client.Close()
}
