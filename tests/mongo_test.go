package tests

import (
	"os"
	"testing"

	"github.com/developia-II/ecommerce-backend/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestConnectToDB_InvalidURI(t *testing.T) {
	os.Setenv("MONGO_URI", "invalid_uri")
	_, err := database.ConnectToDB()
	assert.Error(t, err)
}
