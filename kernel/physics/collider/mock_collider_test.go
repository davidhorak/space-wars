package collider

import (
	"github.com/davidhorak/space-wars/kernel/physics"
	"github.com/stretchr/testify/mock"
)

type MockCollider struct {
	mock.Mock
}

func (m *MockCollider) Enabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockCollider) SetEnabled(enabled bool) {
	m.On("SetEnabled", mock.Anything).Return(nil).Once()
}

func (m *MockCollider) Position() physics.Vector2 {
	args := m.Called()
	return args.Get(0).(physics.Vector2)
}

func (m *MockCollider) SetPosition(position physics.Vector2) {
	m.On("SetPosition", mock.Anything).Return(nil).Once()
}

func (m *MockCollider) Rotation() float64 {
	args := m.Called()
	return args.Get(0).(float64)
}

func (m *MockCollider) SetRotation(rotation float64) {
	m.On("SetRotation", mock.Anything).Return(nil).Once()
}

func (m *MockCollider) CollidesWith(other Collider) bool {
	args := m.Called(other)
	return args.Bool(0)
}

func (m *MockCollider) Serialize() map[string]interface{} {
	args := m.Called()
	return args.Get(0).(map[string]interface{})
}
