package main

import (
	vidio "github.com/AlexEidt/Vidio"
	"github.com/marcoshuck/quic-demo/camera"
	"github.com/marcoshuck/quic-demo/server"
	"log/slog"
	"net"
	"os"
)

const cameraID = 2

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	logger.Debug("Listening on UDP port")
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: 3030,
	})
	if err != nil {
		logger.Error("Failed to listen on udp port 3030", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Debug("Starting camera")
	cam, err := vidio.NewCamera(cameraID)
	if err != nil {
		logger.Error("Failed to initialize camera", slog.Any("error", err))
		os.Exit(2)
	}
	src := camera.NewCamera(cam)

	logger.Debug("Setting up TLS config")
	tlsConfig, err := server.NewTLSConfig()
	if err != nil {
		logger.Error("Failed to initialize TLS config", slog.Any("error", err))
		os.Exit(3)
	}

	logger.Debug("Initializing QUIC server")
	srv := server.NewQuicServer(conn, src, tlsConfig, logger)

	logger.Info("Listening on QUIC Server")
	if err := srv.Listen(); err != nil {
		logger.Error("Failed to run quic server", slog.Any("error", err))
		os.Exit(4)
	}
}
