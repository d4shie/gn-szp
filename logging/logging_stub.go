//go:build !logging_debug

package logging

type Logger struct{}

func (*Logger) Debug(string, ...interface{}) {}
func (*Logger) Info(string, ...interface{})  {}
func (*Logger) Warn(string, ...interface{})  {}
func (*Logger) Error(string, ...interface{}) {}
func (*Logger) Fatal(string, ...interface{}) {}
