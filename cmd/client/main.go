package main

import (
	"context"
	"github.com/marcoshuck/quic-demo/server"
	"github.com/quic-go/quic-go"
	"io"
	"log/slog"
	"os"
	"time"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	logger.Debug("Creating context")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Debug("Creating TLS config")
	tlsConfig, err := server.NewTLSConfig()
	if err != nil {
		logger.Error("Error while creating new TLS config", slog.Any("error", err))
		os.Exit(1)
	}

	tlsConfig.InsecureSkipVerify = true

	logger.Debug("Dialing server on port 3030")
	conn, err := quic.DialAddr(ctx, ":3030", tlsConfig, nil)
	if err != nil {
		logger.Error("Error while dialing", slog.Any("error", err))
		os.Exit(1)
	}
	stream, err := conn.AcceptStream(ctx)
	if err != nil {
		logger.Error("Error while opening stream", slog.Any("error", err))
		os.Exit(1)
	}
	defer stream.Close()

	_, err = io.Copy(os.Stdout, stream)
	if err != nil {
		logger.Error("Error while copying stream into stdout", slog.Any("error", err))
		os.Exit(1)
	}
}
