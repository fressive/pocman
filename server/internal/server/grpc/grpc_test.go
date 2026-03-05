package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/fressive/pocman/server/internal/conf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestTokenAuthInterceptor_NoMetadata(t *testing.T) {
	originalToken := conf.ServerConfig.Server.GRPCToken
	conf.ServerConfig.Server.GRPCToken = "expected-token"
	t.Cleanup(func() {
		conf.ServerConfig.Server.GRPCToken = originalToken
	})

	_, err := tokenAuthInterceptor(
		context.Background(),
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/v1.AgentService/Test"},
		func(ctx context.Context, req any) (any, error) {
			return "ok", nil
		},
	)

	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("expected Unauthenticated, got %v", status.Code(err))
	}

	if status.Convert(err).Message() != "no metadata" {
		t.Fatalf("expected message %q, got %q", "no metadata", status.Convert(err).Message())
	}
}

func TestTokenAuthInterceptor_InvalidToken(t *testing.T) {
	originalToken := conf.ServerConfig.Server.GRPCToken
	conf.ServerConfig.Server.GRPCToken = "expected-token"
	t.Cleanup(func() {
		conf.ServerConfig.Server.GRPCToken = originalToken
	})

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "wrong-token"))

	_, err := tokenAuthInterceptor(
		ctx,
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/v1.AgentService/Test"},
		func(ctx context.Context, req any) (any, error) {
			return "ok", nil
		},
	)

	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("expected Unauthenticated, got %v", status.Code(err))
	}

	if status.Convert(err).Message() != "invalid token" {
		t.Fatalf("expected message %q, got %q", "invalid token", status.Convert(err).Message())
	}
}

func TestTokenAuthInterceptor_EmptyAuthorizationValue(t *testing.T) {
	originalToken := conf.ServerConfig.Server.GRPCToken
	conf.ServerConfig.Server.GRPCToken = "expected-token"
	t.Cleanup(func() {
		conf.ServerConfig.Server.GRPCToken = originalToken
	})

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", ""))

	_, err := tokenAuthInterceptor(
		ctx,
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/v1.AgentService/Test"},
		func(ctx context.Context, req any) (any, error) {
			return "ok", nil
		},
	)

	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("expected Unauthenticated, got %v", status.Code(err))
	}

	if status.Convert(err).Message() != "invalid token" {
		t.Fatalf("expected message %q, got %q", "invalid token", status.Convert(err).Message())
	}
}

func TestTokenAuthInterceptor_ValidTokenCallsHandler(t *testing.T) {
	originalToken := conf.ServerConfig.Server.GRPCToken
	conf.ServerConfig.Server.GRPCToken = "expected-token"
	t.Cleanup(func() {
		conf.ServerConfig.Server.GRPCToken = originalToken
	})

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "expected-token"))

	handlerCalled := false
	expectedErr := errors.New("handler error")

	res, err := tokenAuthInterceptor(
		ctx,
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/v1.AgentService/Test"},
		func(ctx context.Context, req any) (any, error) {
			handlerCalled = true
			return "ok", expectedErr
		},
	)

	if !handlerCalled {
		t.Fatalf("expected handler to be called")
	}

	if res != "ok" {
		t.Fatalf("expected response %q, got %v", "ok", res)
	}

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected handler error to be returned")
	}
}
