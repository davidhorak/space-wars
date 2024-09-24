package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUUID(t *testing.T) {
	assert.Equal(t, int64(1), NewUUID())
	assert.Equal(t, int64(2), NewUUID())
}
