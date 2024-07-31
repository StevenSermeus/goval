package networking

import (
	"crypto/sha256"
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
	isAuth := serverConfig.NoAuth
	for !isAuth {
		command, err := client.Receive()
		if err != nil {
			break
		}
		if command.Command[:4] == "AUTH" {
			passphrase := command.Command[5:]
			h := sha256.New()
			h.Write([]byte(passphrase))
			passphrase = string(h.Sum(nil))
			if passphrase == serverConfig.Passphrase {
				client.Send(types.ResponseInfo{ValueType: "string", Value: "OK"})
				isAuth = true
			} else {
				client.Send(types.ResponseInfo{ValueType: "error", Value: "ERR invalid password"})
			}
		} else {
			client.Send(types.ResponseInfo{ValueType: "error", Value: "ERR Not authenticated"})
		}
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
