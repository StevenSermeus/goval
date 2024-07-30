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
		go commands.ExecuteCommad(&client, command)
	}

	client.Close()
}
