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
}

func (q *quicServer) Listen() error {
	listener, err := q.transport.Listen(q.tlsConfig, q.cfg)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			return err
		}
		go q.handleConn(conn)
	}
}

func Run(srv Server) error {
	return srv.Listen()
}

func (q *quicServer) handleConn(conn quic.Connection) {
	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		return
	}
	defer stream.Close()

	_, err = io.Copy(stream, q.source)
	if err != nil {
		return
	}
}

func NewQuicServer(conn net.PacketConn, source io.Reader, tlsConfig *tls.Config) Server {
	return &quicServer{
		transport: &quic.Transport{Conn: conn},
		source:    source,
		tlsConfig: tlsConfig,
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
