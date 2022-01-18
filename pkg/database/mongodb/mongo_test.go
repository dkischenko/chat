package mongodb

import (
	"context"
	"github.com/dkischenko/chat/internal/config"
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	configPath := os.Getenv("CONFIG")
	cfg := config.GetConfig(configPath)
	storage := cfg.Storage
	_, err := NewClient(context.Background(), storage.Host, storage.Port, storage.Username,
		storage.Password, storage.Database, storage.Options.AuthDB)
	if err != nil {
		t.Errorf("Can not connect to Postgres database due error: %s", err)
	}
}
