package mocks

// MockLogger объект логгера, который необходим для моков в тестах.
type MockLogger struct{}

func (m MockLogger) Debug(args ...interface{}) {}

func (m MockLogger) Info(args ...interface{}) {}

func (m MockLogger) Warn(args ...interface{}) {}

func (m MockLogger) Error(args ...interface{}) {}

func (m MockLogger) Debugf(template string, args ...interface{}) {}

func (m MockLogger) Infof(template string, args ...interface{}) {}

func (m MockLogger) Warnf(template string, args ...interface{}) {}

func (m MockLogger) Errorf(template string, args ...interface{}) {}

func (m MockLogger) Debugw(msg string, keysAndValues ...interface{}) {}

func (m MockLogger) Infow(msg string, keysAndValues ...interface{}) {}

func (m MockLogger) Warnw(msg string, keysAndValues ...interface{}) {}

func (m MockLogger) Errorw(msg string, keysAndValues ...interface{}) {}

func (m MockLogger) Print(args ...interface{}) {}

func (m MockLogger) Printf(format string, args ...interface{}) {}
