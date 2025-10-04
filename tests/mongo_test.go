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

func TestConnectToDB_ValidURI(t *testing.T) {
	os.Setenv("MONGO_URI", "mongodb+srv://vendora:jebCkbDHuCgbvyCy@cluster0.um23o.mongodb.net/")
	_, err := database.ConnectToDB()
	assert.Error(t, err)
}
