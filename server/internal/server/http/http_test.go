package http

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fressive/pocman/common/pkg/model"
	"github.com/fressive/pocman/server/internal/conf"
	"github.com/fressive/pocman/server/internal/data"
	"github.com/fressive/pocman/server/internal/model/dto"
	"github.com/fressive/pocman/server/internal/util"
	"github.com/gin-gonic/gin"
)

func setupHTTPTokenTestDB(t *testing.T) {
	t.Helper()

	originalDriver := conf.ServerConfig.Data.Database.Driver
	originalSource := conf.ServerConfig.Data.Database.Source

	conf.ServerConfig.Data.Database.Driver = "sqlite"
	conf.ServerConfig.Data.Database.Source = "file::memory:?cache=shared"

	t.Cleanup(func() {
		conf.ServerConfig.Data.Database.Driver = originalDriver
		conf.ServerConfig.Data.Database.Source = originalSource
	})

	if err := data.InitDatabase(); err != nil {
		t.Fatalf("failed to init test database: %v", err)
	}

	if err := data.DB.AutoMigrate(&dto.APIToken{}); err != nil {
		t.Fatalf("failed to migrate APIToken: %v", err)
	}
}

func TestTokenAuthMiddleware_NoAuthorizationHeader(t *testing.T) {
	setupHTTPTokenTestDB(t)
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(TokenAuthMiddleware())
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}

	var body model.Response[map[string]any]
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body.Code != 20000 {
		t.Fatalf("expected response code 20000, got %d", body.Code)
	}
}

func TestTokenAuthMiddleware_InvalidBase64Token(t *testing.T) {
	setupHTTPTokenTestDB(t)
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(TokenAuthMiddleware())
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer !!!")
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var body model.Response[map[string]any]
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body.Code != 20001 {
		t.Fatalf("expected response code 20001, got %d", body.Code)
	}
}

func TestTokenAuthMiddleware_MalformedAuthorizationHeader(t *testing.T) {
	setupHTTPTokenTestDB(t)
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(TokenAuthMiddleware())
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer")
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}

	var body model.Response[map[string]any]
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body.Code != 20000 {
		t.Fatalf("expected response code 20000, got %d", body.Code)
	}
}

func TestTokenAuthMiddleware_InvalidToken(t *testing.T) {
	setupHTTPTokenTestDB(t)
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(TokenAuthMiddleware())
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+base64.RawURLEncoding.EncodeToString([]byte("missing-token")))
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}

	var body model.Response[map[string]any]
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body.Code != 20002 {
		t.Fatalf("expected response code 20002, got %d", body.Code)
	}
}

func TestTokenAuthMiddleware_ValidTokenAllowsRequest(t *testing.T) {
	setupHTTPTokenTestDB(t)
	gin.SetMode(gin.TestMode)

	token := "valid-token"
	if err := data.DB.Create(&dto.APIToken{
		Hash:        util.HashToken(token),
		Name:        "test",
		ValidBefore: time.Now().Add(30 * time.Minute),
		LastUsed:    time.Now(),
	}).Error; err != nil {
		t.Fatalf("failed to seed api token: %v", err)
	}

	r := gin.New()
	r.Use(TokenAuthMiddleware())
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+base64.RawURLEncoding.EncodeToString([]byte(token)))
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestNewHTTPServer_PingRoute(t *testing.T) {
	setupHTTPTokenTestDB(t)
	gin.SetMode(gin.TestMode)

	token := "ping-token"
	if err := data.DB.Create(&dto.APIToken{
		Hash:        util.HashToken(token),
		Name:        "ping",
		ValidBefore: time.Now().Add(30 * time.Minute),
		LastUsed:    time.Now(),
	}).Error; err != nil {
		t.Fatalf("failed to seed api token: %v", err)
	}

	r, err := NewHTTPServer(NewPingHandler(), NewAgentHandler(), NewFileHandler(), NewVulnHandler())
	if err != nil {
		t.Fatalf("failed to build http server: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
	req.Header.Set("Authorization", "Bearer "+base64.RawURLEncoding.EncodeToString([]byte(token)))
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var body model.Response[map[string]string]
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body.Code != 0 {
		t.Fatalf("expected response code 0, got %d", body.Code)
	}

	if body.Data["ping"] != "pong" {
		t.Fatalf("expected ping response to be pong, got %q", body.Data["ping"])
	}
}
