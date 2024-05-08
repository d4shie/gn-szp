//go:build logging_debug

package logging

import "log"

type Logger struct{}

func (*Logger) Debug(format string, v ...interface{}) {
	log.Printf("[DEBUG]: "+format, v...)
}
func (*Logger) Info(format string, v ...interface{}) {
	log.Printf("[INFO]: "+format, v...)
}
func (*Logger) Warn(format string, v ...interface{}) {
	log.Printf("[WARN]: "+format, v...)
}
func (*Logger) Error(format string, v ...interface{}) {
	log.Printf("[ERROR]: "+format, v...)
}
func (*Logger) Fatal(format string, v ...interface{}) {
	log.Fatalf("[FATAL]: "+format, v...)
}
