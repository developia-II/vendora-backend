package tests

import (
	"os"
	"testing"
	"time"

	"github.com/developia-II/ecommerce-backend/utils"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key-12345")
	defer os.Unsetenv("JWT_SECRET")
	userId := "0283721"
	userRole := "vendor"
	duration := time.Duration(24 * time.Hour)
	token, err := utils.GenerateToken(userId, userRole, duration)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
