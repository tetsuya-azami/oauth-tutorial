package test

type MockLogger struct {
}

func NewMockLogger() *MockLogger {
	return &MockLogger{}
}

func (m *MockLogger) Info(msg string, args ...any) {
}

func (m *MockLogger) Error(msg string, args ...any) {
}

func (m *MockLogger) Warn(msg string, args ...any) {
}
