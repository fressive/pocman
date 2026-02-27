package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"
	"os"

	"github.com/fressive/pocman/agent/internal/conf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type TokenAuth struct {
	Token string
}

func (t TokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": t.Token,
	}, nil
}

func (t TokenAuth) RequireTransportSecurity() bool {
	return false
}

func NewConn() (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	var err error

	server := conf.AgentConfig.Server

	if server == nil {
		return nil, fmt.Errorf("configuration incomplete: missing server section")
	}

	if server.Host == nil || *server.Host == "" {
		return nil, fmt.Errorf("configuration incomplete: missing host field")
	}

	if server.Port == nil {
		return nil, fmt.Errorf("configuration incomplete: missing port field")
	}

	target := fmt.Sprintf("%s:%d", *server.Host, *server.Port)

	if grpcCert := server.GRPCCert; grpcCert != nil {
		// use cert
		cert, err := tls.LoadX509KeyPair(grpcCert.Cert, grpcCert.Key)
		if err != nil {
			return nil, err
		}

		certPool := x509.NewCertPool()
		ca, err := os.ReadFile(grpcCert.CA)
		if err != nil {
			return nil, err
		}
		certPool.AppendCertsFromPEM(ca)

		creds := credentials.NewTLS(&tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{cert},
			RootCAs:      certPool,
		})

		conn, err = grpc.NewClient(target,
			grpc.WithTransportCredentials(creds))
	} else if server.GRPCToken != nil && *server.GRPCToken != "" {
		// use token
		slog.Warn("Transport will be insecure when token is used as the authenticator. " +
			"Please configure certificates as an alternative method.")

		auth := TokenAuth{Token: *server.GRPCToken}
		conn, err = grpc.NewClient(target,
			grpc.WithPerRPCCredentials(auth),
			grpc.WithInsecure())
	} else {
		return nil, fmt.Errorf("Authentication method not configured (cert|token)")
	}

	return conn, err
}
