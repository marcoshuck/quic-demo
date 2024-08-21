package server

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"github.com/quic-go/quic-go"
	"io"
	"log/slog"
	"math/big"
	"net"
)

type Server interface {
	Listen() error
}

type quicServer struct {
	tlsConfig *tls.Config
	cfg       *quic.Config
	transport *quic.Transport
	source    io.Reader
	logger    *slog.Logger
}

func (q *quicServer) Listen() error {
	q.logger.Debug("Listening...")
	listener, err := q.transport.Listen(q.tlsConfig, q.cfg)
	if err != nil {
		return err
	}
	defer listener.Close()
	q.logger.Debug("Waiting on a new listener...")
	conn, err := listener.Accept(context.Background())
	if err != nil {
		return err
	}
	q.logger.Debug("Handling incoming connection...")
	err = q.handleConn(conn)
	if err != nil {
		return err
	}
	return nil
}

func Run(srv Server) error {
	return srv.Listen()
}

func (q *quicServer) handleConn(conn quic.Connection) error {
	q.logger.Debug("Accepting stream...")
	stream, err := conn.OpenStream()
	if err != nil {
		return err
	}
	defer stream.Close()

	q.logger.Debug("Copying camera data into stream...")
	_, err = io.Copy(stream, q.source)
	if err != nil {
		return err
	}
	return nil
}

func NewQuicServer(conn net.PacketConn, source io.Reader, tlsConfig *tls.Config, logger *slog.Logger) Server {
	return &quicServer{
		transport: &quic.Transport{Conn: conn},
		source:    source,
		tlsConfig: tlsConfig,
		logger:    logger,
	}
}

func NewTLSConfig() (*tls.Config, error) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, err
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"camera-quic-server"},
	}, nil
}
