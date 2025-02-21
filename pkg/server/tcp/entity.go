package tcp

import (
	"github.com/skolldire/web-simplify/pkg/utilities/log"
)

type Service interface {
	GetMessage(f ProcessingFunc)
}

type ProcessingFunc func(msg string) (string, error)

type Dependencies struct {
	Config Config
	Log    log.Service
}

type Config struct {
	Port         string `json:"port"`
	InstanceName string `json:"instance_name"`
}
