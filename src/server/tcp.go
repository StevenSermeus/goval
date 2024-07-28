package server

import (
	"errors"
	"net"
	"strconv"
	"strings"

	"github.com/StevenSermeus/goval/src/cache"
	"github.com/StevenSermeus/goval/src/command"
	"github.com/StevenSermeus/goval/src/config"
	"github.com/StevenSermeus/goval/src/logging"
	"golang.org/x/net/netutil"
)

func TcpServer(cache *cache.Cache, serverConfig *config.Config) {
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
			continue
		}
		go handleClient(conn, cache, serverConfig)
	}
}

func handleClient(conn net.Conn, cache *cache.Cache, serverConfig *config.Config) {
	logging.Info.Println("Client connected from", conn.RemoteAddr())
	defer conn.Close()
	for {
		logging.Info.Println("Waiting for command from client")
		commandString, err := waitForCommand(conn, cache, serverConfig)
		if err != nil {
			if _, ok := err.(EOCError); ok {
				logging.Info.Println("Client disconnected from", conn.RemoteAddr())
				break
			}
			logging.Error.Println("Error reading command from client:", err)
			continue
		}

		response, err := command.Exec(commandString, cache, serverConfig)
		if err != nil {
			logging.Error.Println("Error executing command:", err)
			conn.Write([]byte("ERROR: " + err.Error()))
			continue
		}
		if response != "" {
			conn.Write([]byte(response))
		}else{
			conn.Write([]byte("OK"))
		}
	}
	logging.Info.Println("Client disconnected from", conn.RemoteAddr())
}

func collectAllBuffer(conn net.Conn, n int, data string) (string,error) {
	for i := 0; i < n; i++ {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			return "", EOCError{}
		}
		if n == 0 {
			return "", EOCError{}
		}
		data += string(buffer[:n])
	}
	return data, nil
}

func waitForCommand(conn net.Conn, cache *cache.Cache, serverConfig *config.Config) (string, error){
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil || n == 0 {
		return "", EOCError{}
	}	
	netString := string(buffer[:n])
	netStringSplit := strings.Split(netString, ":")
	if len(netStringSplit) < 2 {
		return "", errors.New("invalid command format")
	}
	commandString := strings.Join(netStringSplit[1:], "")
	messageSize, err := strconv.Atoi(netStringSplit[0])
	if err != nil {
		return "", err
	}
	if messageSize / serverConfig.BufferSize > 0 {
		commandString, err = collectAllBuffer(conn, messageSize / serverConfig.BufferSize, commandString)
		if err != nil {
			return "", err
		}
	}
	return commandString, nil
}

type EOCError struct {
}

func (e EOCError) Error() string {
	return "End of connection"
}

