package fSchedule

import "strings"

type healthCheck struct {
}

func (c *healthCheck) Check() (string, error) {
	err := defaultClient.RegistryClient()
	return "FSchedule." + strings.Join(defaultServer.Address, ","), err
}
