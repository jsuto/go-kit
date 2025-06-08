package logx

import (
    "log"
    "os"
    "strings"
)

const (
    LOG_ERROR = iota
    LOG_WARN
    LOG_INFO
    LOG_DEBUG
)

var currentLevel = LOG_INFO

func SetLogLevel() {
    levelStr := strings.ToUpper(os.Getenv("LOG_LEVEL"))

    switch levelStr {
    case "ERROR":
        currentLevel = LOG_ERROR
    case "WARN":
        currentLevel = LOG_WARN
    case "INFO":
        currentLevel = LOG_INFO
    case "DEBUG":
        currentLevel = LOG_DEBUG
    case "":
        // No env var set; default silently
        currentLevel = LOG_INFO
    default:
        currentLevel = LOG_INFO
        log.Printf("[WARN] Invalid LOG_LEVEL '%s'; defaulting to LOG_INFO", levelStr)
    }
}

func logf(level int, prefix string, format string, args ...interface{}) {
    if level <= currentLevel {
        log.Printf("[" + prefix + "] " + format, args...)
    }
}

func Error(format string, args ...interface{}) { logf(LOG_ERROR, "ERROR", format, args...) }
func Warn(format string, args ...interface{})  { logf(LOG_WARN, "WARN", format, args...) }
func Info(format string, args ...interface{})  { logf(LOG_INFO, "INFO", format, args...) }
func Debug(format string, args ...interface{}) { logf(LOG_DEBUG, "DEBUG", format, args...) }
