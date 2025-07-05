package api

import (
	"os"
	db "simple_bank/db/sqlc"
	"simple_bank/util"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func NewTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	if err != nil {
		t.Fatalf("cannot create test server: %v", err)
	}

	return server

}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
