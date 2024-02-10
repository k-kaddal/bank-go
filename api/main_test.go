package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/k-kaddal/bank-go/db/sqlc"
	"github.com/k-kaddal/bank-go/util"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func newTestServer (t *testing.T, store db.Store) *Server {
	config := util.Config{
		Token_Symmetric_Key: util.RandomString(32),
		Access_Token_Duration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)
	
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}