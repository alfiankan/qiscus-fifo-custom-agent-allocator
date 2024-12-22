package utils

import (
	"log"
	"strings"
)

const (
	LOG_INFO  = "INFO"
	LOG_DEBUG = "DEBUG"
	LOG_ERROR = "ERROR"
)

func LogWrite(log_name string, level string, messages ...string) {
	log.Printf("[%s][%s] %s", log_name, level, strings.Join(messages, " - "))
}
