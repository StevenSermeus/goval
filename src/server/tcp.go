package server

import (
	"net"
	"time"

	"github.com/StevenSermeus/goval/src/cache"
	"github.com/StevenSermeus/goval/src/command"
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
			continue
		}
		go handleClient(conn, cache, serverConfig)
	}
}

func handleClient(conn net.Conn, cache *cache.Cache, serverConfig *config.Config) {
	logging.Info.Println("Client connected from", conn.RemoteAddr())
	defer conn.Close()
	for {
		message, err := collectAllBuffer(conn)
		startExecution := time.Now()
		if err != nil {
			if _, ok := err.(EOCError); ok {
				logging.Info.Println("Client disconnected from", conn.RemoteAddr())
				break
			}
			logging.Error.Println("Error reading command from client:", err)
		}
		if message[0] != '!' {
			logging.Error.Println("Invalid command format")
			conn.Write([]byte("ERROR: invalid command format\n\r\n\r"))
			continue
		}
		//remove the ! from the start and the \n\r\n\r from the end
		response, err := command.Exec(message[1:len(message)-4], cache, serverConfig)
		if err != nil {
			logging.Info.Println("Failed to execute", time.Since(startExecution))
			logging.Error.Println("Error executing command:", err)
			res := "-ERROR: " + err.Error()
			conn.Write([]byte(res + "\n\r\n\r"))
			continue
		}
		if response != "" {
			conn.Write([]byte("+" + response + "\n\r\n\r"))
		} else {
			conn.Write([]byte("+OK\n\r\n\r"))
		}
		logging.Info.Println("Executed command in", time.Since(startExecution))
	}
	logging.Info.Println("Client disconnected from", conn.RemoteAddr())
}

func collectAllBuffer(conn net.Conn) (string, error) {
	message := ""
	for {
		buffer := make([]byte, 1024)
		n := 0
		var err error
		n, err =
			conn.Read(buffer)
		if err != nil {
			return "", EOCError{}
		}
		if n == 0 {
			return "", EOCError{}
		}
		message += string(buffer[:n])
		if message[len(message)-4:] == "\n\r\n\r" {
			break
		}
	}
	return message, nil
}

// func collectAllBuffer(conn net.Conn, n int, data string) (string, error) {
// 	for i := 0; i < n; i++ {
// 		buffer := make([]byte, 1024)
// 		n, err := conn.Read(buffer)
// 		if err != nil {
// 			return "", EOCError{}
// 		}
// 		if n == 0 {
// 			return "", EOCError{}
// 		}
// 		data += string(buffer[:n])
// 	}
// 	return data, nil
// }

// func waitForCommand(conn net.Conn, serverConfig *config.Config) (string, error) {
// 	buffer := make([]byte, 1024)
// 	n, err := conn.Read(buffer)
// 	if err != nil || n == 0 {
// 		return "", EOCError{}
// 	}
// 	netString := string(buffer[:n])
// 	netStringSplit := strings.Split(netString, ":")
// 	if len(netStringSplit) < 2 {
// 		return "", errors.New("invalid command format")
// 	}
// 	commandString := strings.Join(netStringSplit[1:], "")
// 	messageSize, err := strconv.Atoi(netStringSplit[0])
// 	if err != nil {
// 		return "", err
// 	}
// 	if messageSize/serverConfig.BufferSize > 0 {
// 		commandString, err = collectAllBuffer(conn, messageSize/serverConfig.BufferSize, commandString)
// 		if err != nil {
// 			return "", err
// 		}
// 	}
// 	return commandString, nil
// }

type EOCError struct {
}

func (e EOCError) Error() string {
	return "End of connection"
}

// isAuthenticated := false
// for !isAuthenticated {
// 	buffer := make([]byte, 1024)
// 	n, err :=
// 		conn.Read(buffer)
// 	if err != nil || n == 0 {
// 		logging.Error.Println("Error reading from client:", err)
// 		return
// 	}
// 	commandString := string(buffer[:n])
// 	if strings.TrimSpace(commandString) == serverConfig.Passphrase {
// 		conn.Write([]byte("OK"))
// 		isAuthenticated = true
// 	} else {
// 		conn.Write([]byte("ERROR: invalid password"))
// 	}
// }
// logging.Info.Println("Client authenticated")
// for {
// 	logging.Info.Println("Waiting for command from client")
// 	commandString, err := waitForCommand(conn, serverConfig)
// 	if err != nil {
// 		if _, ok := err.(EOCError); ok {
// 			logging.Info.Println("Client disconnected from", conn.RemoteAddr())
// 			break
// 		}
// 		logging.Error.Println("Error reading command from client:", err)
// 		continue
// 	}
// 	startExecution := time.Now()
// 	response, err := command.Exec(commandString, cache, serverConfig)
// 	if err != nil {
// 		logging.Error.Println("Error executing command:", err)
// 		res := "ERROR: " + err.Error()
// 		conn.Write([]byte(string(rune(len(res))) + ":" + res))
// 		continue
// 	}
// 	logging.Info.Println("Command executed in", time.Since(startExecution))
// 	if response != "" {
// 		conn.Write([]byte(string(rune(len(response))) + ":" + response))
// 	} else {
// 		conn.Write([]byte("2:OK"))
// 	}

// }
