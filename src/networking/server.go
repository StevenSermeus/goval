package networking

import (
	"net"

	"github.com/StevenSermeus/goval/src/cache"
	"github.com/StevenSermeus/goval/src/config"
	"github.com/StevenSermeus/goval/src/logging"
	"golang.org/x/net/netutil"
)

func Tcp(cache *cache.Cache, serverConfig *config.Config) {
	listener, err := net.Listen("tcp", "0.0.0.0:"+serverConfig.Port)
	listener = netutil.LimitListener(listener, serverConfig.MaxConnections)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	logging.Info.Println("Listening on port", serverConfig.Port, "with max connections", serverConfig.MaxConnections)
	for {
		conn, err := listener.Accept()
		if err != nil {
			logging.Error.Println("Error accepting connection", err)
			continue
		}
		go HandleClient(conn, cache, serverConfig)
	}
}
