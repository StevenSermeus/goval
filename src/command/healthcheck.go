package command

import (
	"github.com/StevenSermeus/goval/src/logging"
)

func HealthCheck() (string, error) {
	logging.Info.Println("Executing healthcheck")
	return "OK", nil
}

func IsHealthCheckCommand(command string) bool {
	return command == "HEALTHCHECK"
}

func ExecuteHealthCheckCommand(command string) (string, error) {
	return HealthCheck()
}