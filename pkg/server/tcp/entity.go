package tcp

import (
	"context"
	"github.com/skolldire/web-simplify/pkg/utilities/log"
	"net"
)

type Service interface {
	Start(ctx context.Context, f ProcessingFunc)
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

type service struct {
	server net.Listener
	log    log.Service
	port   string
}
