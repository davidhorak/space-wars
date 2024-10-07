package game

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMessage_Serialize(t *testing.T) {
	now := time.Now()
	message := Message{
		id:      1,
		logType: LogTypeDamage,
		time:    now,
		message: "test",
		meta:    map[string]interface{}{"test": "test"},
	}

	serialized := message.Serialize()
	assert.Equal(t, int64(1), serialized["id"])
	assert.Equal(t, "damage", serialized["logType"])
	assert.Equal(t, now.Format("2006-01-02 15:04:05"), serialized["time"])
	assert.Equal(t, "test", serialized["message"])
	assert.Equal(t, map[string]interface{}{"test": "test"}, serialized["meta"])
}

func TestNewLogger(t *testing.T) {
	logger := NewLogger()
	assert.Equal(t, []Message{}, logger.Logs())
}

func TestLogger_Logs(t *testing.T) {
	logger := NewLogger()
	assert.Equal(t, []Message{}, logger.Logs())
}

func TestLogger_Clear(t *testing.T) {
	logger := NewLogger()
	logger.Clear()
	assert.Equal(t, []Message{}, logger.Logs())
}

func TestLogger_AddMessage(t *testing.T) {
	logger := NewLogger()
	logger.AddMessage(Message{
		id:      1,
		logType: LogTypeDamage,
		time:    time.Now(),
		message: "test",
		meta:    map[string]interface{}{"test": "test"},
	})

	assert.Equal(t, 1, len(logger.Logs()))
	assert.Equal(t, Message{
		id:      1,
		logType: LogTypeDamage,
		time:    time.Now(),
		message: "test",
		meta:    map[string]interface{}{"test": "test"},
	}, logger.Logs()[0])
}

func TestLogger_Damage(t *testing.T) {
	now := time.Now()
	logger := NewLogger()
	logger.Damage(now, 10, "test", "other", DamageTypeUnknown)

	log := logger.Logs()[0]
	assert.Equal(t, 1, len(logger.Logs()))
	assert.Greater(t, log.id, int64(0))
	assert.Equal(t, LogTypeDamage, log.logType)
	assert.Equal(t, now, log.time)
	assert.Equal(t, "\"test\" did 10.00 damage to \"other\" with unknown", log.message)
	assert.Equal(t, map[string]interface{}{"who": "test", "whom": "other", "damage": "10.00", "damageType": "unknown"}, log.meta)
}

func TestLogger_Kill(t *testing.T) {
	now := time.Now()
	logger := NewLogger()
	logger.Kill(now, "test", "test")

	log := logger.Logs()[0]
	assert.Equal(t, 1, len(logger.Logs()))
	assert.Greater(t, log.id, int64(0))
	assert.Equal(t, LogTypeKill, log.logType)
	assert.Equal(t, now, log.time)
	assert.Equal(t, "\"test\" was killed by \"test\"", log.message)
	assert.Equal(t, map[string]interface{}{"who": "test", "whom": "test"}, log.meta)
}

func TestLogger_Collision(t *testing.T) {
	now := time.Now()
	logger := NewLogger()
	logger.Collision(now, "test", "test")

	log := logger.Logs()[0]
	assert.Equal(t, 1, len(logger.Logs()))
	assert.Greater(t, log.id, int64(0))
	assert.Equal(t, LogTypeCollision, log.logType)
	assert.Equal(t, now, log.time)
	assert.Equal(t, "\"test\" collided with \"test\"", log.message)
	assert.Equal(t, map[string]interface{}{"who": "test", "with": "test"}, log.meta)
}

func TestLogger_GameState(t *testing.T) {
	now := time.Now()
	logger := NewLogger()
	logger.GameState(now, Running)

	log := logger.Logs()[0]
	assert.Equal(t, 1, len(logger.Logs()))
	assert.Greater(t, log.id, int64(0))
	assert.Equal(t, LogTypeGameState, log.logType)
	assert.Equal(t, now, log.time)
	assert.Equal(t, "Game state changed to: running", log.message)
	assert.Equal(t, map[string]interface{}{"state": "running"}, log.meta)
}
