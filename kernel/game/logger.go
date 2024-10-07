package game

import (
	"fmt"
	"time"
)

type LogType string

const (
	LogTypeDamage    LogType = "damage"
	LogTypeKill      LogType = "kill"
	LogTypeCollision LogType = "collision"
	LogTypeGameState LogType = "game_state"
)

type Message struct {
	id      int64
	logType LogType
	time    time.Time
	message string
	meta    map[string]interface{}
}

func (message *Message) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"id":      message.id,
		"logType": string(message.logType),
		"time":    message.time.Format("2006-01-02 15:04:05"),
		"message": message.message,
		"meta":    message.meta,
	}
}

type Logger interface {
	Logs() []Message
	Clear()
	AddMessage(message Message)
	Damage(time time.Time, damage float64, who string, by string, damageType DamageType)
	Kill(time time.Time, who string, by string)
	Collision(time time.Time, who string, with string)
	GameState(time time.Time, state Status)
}

func NewLogger() Logger {
	return &logger{
		messages: []Message{},
	}
}

type logger struct {
	messages []Message
}

func (logger *logger) Logs() []Message {
	return logger.messages
}

func (logger *logger) Clear() {
	logger.messages = []Message{}
}

func (logger *logger) AddMessage(message Message) {
	logger.messages = append(logger.messages, message)
}

func (logger *logger) Damage(time time.Time, damage float64, who string, whom string, damageType DamageType) {
	logger.messages = append(logger.messages, Message{
		id:      NewUUID(),
		logType: LogTypeDamage,
		time:    time,
		message: fmt.Sprintf("\"%s\" did %.2f damage to \"%s\" with %s", who, damage, whom, damageType),
		meta: map[string]interface{}{
			"who":        who,
			"whom":       whom,
			"damage":     fmt.Sprintf("%.2f", damage),
			"damageType": string(damageType),
		},
	})
}

func (logger *logger) Kill(time time.Time, who string, whom string) {
	logger.messages = append(logger.messages, Message{
		id:      NewUUID(),
		logType: LogTypeKill,
		time:    time,
		message: fmt.Sprintf("\"%s\" was killed by \"%s\"", who, whom),
		meta: map[string]interface{}{
			"who":  who,
			"whom": whom,
		},
	})
}

func (logger *logger) Collision(time time.Time, who string, with string) {
	logger.messages = append(logger.messages, Message{
		id:      NewUUID(),
		logType: LogTypeCollision,
		time:    time,
		message: fmt.Sprintf("\"%s\" collided with \"%s\"", who, with),
		meta: map[string]interface{}{
			"who":  who,
			"with": with,
		},
	})
}

func (logger *logger) GameState(time time.Time, state Status) {
	logger.messages = append(logger.messages, Message{
		id:      NewUUID(),
		logType: LogTypeGameState,
		time:    time,
		message: fmt.Sprintf("Game state changed to: %s", state),
		meta: map[string]interface{}{
			"state": string(state),
		},
	})
}
