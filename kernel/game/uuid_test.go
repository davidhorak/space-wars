package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUUID(t *testing.T) {
	ResetUUID()
	assert.Equal(t, int64(1), NewUUID())
	assert.Equal(t, int64(2), NewUUID())
}

func TestGetUUID(t *testing.T) {
	ResetUUID()
	NewUUID()
	assert.Equal(t, int64(1), GetUUID())
}

func TestResetUUID(t *testing.T) {
	ResetUUID()
	NewUUID()
	NewUUID()
	assert.Equal(t, int64(3), NewUUID())
	ResetUUID()
	assert.Equal(t, int64(1), NewUUID())
}

func TestSetUUID(t *testing.T) {
	ResetUUID()
	SetUUID(100)
	assert.Equal(t, int64(101), NewUUID())
}
