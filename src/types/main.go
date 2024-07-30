package types

import (
	"errors"
	"net"

	"github.com/StevenSermeus/goval/src/cache"
	"github.com/StevenSermeus/goval/src/config"
	"github.com/StevenSermeus/goval/src/logging"
	"github.com/StevenSermeus/goval/src/utils"
)

type EOCError struct {
}

func (e EOCError) Error() string {
	return "End of connection"
}

type Client struct {
	Conn         net.Conn
	Cache        *cache.Cache
	ServerConfig *config.Config
}

type ResponseInfo struct {
	ValueType string
	Value     any
}

func (c *Client) Send(data ResponseInfo) {
	logging.Info.Println("Sending response :", data.Value)
	valueCode, err := utils.GetCodeFromType(data.ValueType)
	if err != nil {
		c.Conn.Write([]byte("-" + err.Error() + "\n\r\n\r"))
	}
	c.Conn.Write([]byte(valueCode + data.Value.(string) + "\n\r\n\r"))
}

func (c *Client) Receive() (CommandInfo, error) {
	message := ""
	for {
		buffer := make([]byte, 1024)
		n := 0
		var err error
		n, err = c.Conn.Read(buffer)
		if err != nil {
			return CommandInfo{}, EOCError{}
		}
		if n == 0 {
			return CommandInfo{}, EOCError{}
		}
		message += string(buffer[:n])
		if message[len(message)-4:] == "\n\r\n\r" {
			break
		}
	}
	return c.parseNetworkCommand(message)
}

func (c *Client) Close() {
	c.Conn.Close()
}

// The Value Type is determined by the first character of the message:
// - '+' indicates a string command
// - '-' indicates an error command
// - ':' indicates an integer command
// - '!' indicates a generic command
type CommandInfo struct {
	Command   string
	ValueType string
}

func (c *Client) parseNetworkCommand(message string) (CommandInfo, error) {
	switch message[0] {
	case '+':
		return CommandInfo{Command: message[1 : len(message)-4], ValueType: "string"}, nil
	case '-':
		return CommandInfo{Command: message[1 : len(message)-4], ValueType: "error"}, nil
	case ':':
		return CommandInfo{Command: message[1 : len(message)-4], ValueType: "int"}, nil
	case '!':
		return CommandInfo{Command: message[1 : len(message)-4], ValueType: "command"}, nil
	default:
		return CommandInfo{}, errors.New("invalid message format")
	}
}
