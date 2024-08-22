package server

import (
	"fmt"
	"go.uber.org/zap"
	"parkingCam/internal/logger"
	"parkingCam/internal/transport"
	"time"
)

const serviceName = "cam-server"

type Config struct {
	RecordsLifetime time.Duration
	Token           string
}

type Server struct {
	logger   *zap.Logger
	receiver *transport.Receiver
}

func (s *Server) ListenAndServe() error {
	s.logger.Info("server start listening")

	return nil
}

func (s *Server) acceptBucket(bucket transport.Bucket) {
	// TODO
}

func Run() error {
	l, err := logger.BuildLogger(serviceName)
	if err != nil {
		return fmt.Errorf("cannot build logger for service: %w", err)
	}

	server := Server{logger: l, receiver: transport.NewReceiver()}
	if err = server.ListenAndServe(); err != nil {
		return fmt.Errorf("cannot start server: %w", err)
	}

	return nil
}
