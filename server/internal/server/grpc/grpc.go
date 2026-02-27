package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"
	"net"
	"os"

	v1 "github.com/fressive/pocman/common/proto/v1"
	"github.com/fressive/pocman/server/internal/conf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	v1.UnimplementedAgentServiceServer
}

func tokenAuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "no metadata")
	}

	if tokens, ok := meta["authorization"]; ok {
		if tokens[0] == conf.ServerConfig.Server.GRPCToken {
			return handler(ctx, req)
		}
	}

	return nil, status.Errorf(codes.Unauthenticated, "invalid token")
}

func NewGRPCServer() (*grpc.Server, error) {
	var svr *grpc.Server

	if grpcCert := conf.ServerConfig.Server.GRPCCert; grpcCert != nil {
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
			ClientCAs:    certPool,
		})

		svr = grpc.NewServer(grpc.Creds(creds))
	} else if token := conf.ServerConfig.Server.GRPCToken; token != "" {
		// use token
		slog.Warn("Using token in producton environment is unsafe and will be deprecated in future versions. " +
			"Please configure certificates as an alternative method.")

		svr = grpc.NewServer(grpc.UnaryInterceptor(tokenAuthInterceptor))
	} else {
		return nil, fmt.Errorf("authentication method is not configured")
	}

	v1.RegisterAgentServiceServer(svr, &GRPCServer{})

	return svr, nil
}

func RunGRPCServer() (*grpc.Server, error) {
	srv, err := NewGRPCServer()
	if err != nil {
		return nil, err
	}

	grpcPort := conf.ServerConfig.Server.GRPCPort
	if grpcPort == 0 {
		grpcPort = conf.ServerConfig.Server.Port + 1
	}
	addr := fmt.Sprintf("%s:%d", conf.ServerConfig.Server.Host, grpcPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	go func() {
		slog.Info("Starting Pocman gRPC server", "addr", addr)
		if err := srv.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			slog.Error("error occured when gRPC serves", "err", err)
		}
	}()

	return srv, nil
}
