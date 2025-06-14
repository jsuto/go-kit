package logx

import (
    "bytes"
    "log"
    "os"
    "strings"
    "testing"
)

func TestSetLogLevelFromEnv(t *testing.T) {
    tests := []struct {
        envValue    string
        expected    int
        shouldWarn  bool
    }{
        {"DEBUG", LOG_DEBUG, false},
        {"INFO", LOG_INFO, false},
        {"WARN", LOG_WARN, false},
        {"ERROR", LOG_ERROR, false},
        {"", LOG_INFO, false},
        {"INVALID_LEVEL", LOG_INFO, true},
    }

    for _, tt := range tests {
        os.Setenv("LOG_LEVEL", tt.envValue)

        var buf bytes.Buffer
        log.SetOutput(&buf)
        SetLogLevel()

        if currentLevel != tt.expected {
            t.Errorf("LOG_LEVEL=%q: expected level %d, got %d", tt.envValue, tt.expected, currentLevel)
        }

        output := buf.String()
        if tt.shouldWarn && !strings.Contains(output, "Invalid LOG_LEVEL") {
            t.Errorf("Expected warning for LOG_LEVEL=%q but got none", tt.envValue)
        }
        if !tt.shouldWarn && strings.Contains(output, "Invalid LOG_LEVEL") {
            t.Errorf("Unexpected warning for LOG_LEVEL=%q", tt.envValue)
        }
    }
}

func TestLoggingLevels(t *testing.T) {
    // Set to INFO for this test
    currentLevel = LOG_INFO

    var buf bytes.Buffer
    log.SetOutput(&buf)

    // Should not log: DEBUG
    Debug("debug log")
    if strings.Contains(buf.String(), "DEBUG") {
        t.Error("Debug log should not be printed at LOG_INFO level")
    }
    buf.Reset()

    // Should log: INFO
    Info("info log")
    if !strings.Contains(buf.String(), "INFO") {
        t.Error("Info log should be printed at LOG_INFO level")
    }
    buf.Reset()

    // Should log: WARN
    Warn("warn log")
    if !strings.Contains(buf.String(), "WARN") {
        t.Error("Warn log should be printed at LOG_INFO level")
    }
    buf.Reset()

    // Should log: ERROR
    Error("error log")
    if !strings.Contains(buf.String(), "ERROR") {
        t.Error("Error log should be printed at LOG_INFO level")
    }
}
