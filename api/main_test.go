package api

import (
	"btc_billionaire/postgres"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func newTestServer(t *testing.T, store postgres.Store) *Handler {
	return NewServer(store)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
