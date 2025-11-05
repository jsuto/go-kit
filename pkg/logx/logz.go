package logx

import (
    "context"
    "fmt"
    "log"
    "os"
    "strings"
    "time"

    "github.com/rs/zerolog"
)

const (
    LOG_ERROR = iota
    LOG_WARN
    LOG_INFO
    LOG_DEBUG
)

var (
    currentLevel = LOG_INFO
    jsonMode     = true
    zl           zerolog.Logger
)

const requestIDKey = "requestid"

func SetLogLevel(logLevel, logFormat string) {
    levelStr := strings.ToUpper(logLevel)
    switch levelStr {
    case "ERROR":
        currentLevel = LOG_ERROR
    case "WARN":
        currentLevel = LOG_WARN
    case "INFO":
        currentLevel = LOG_INFO
    case "DEBUG":
        currentLevel = LOG_DEBUG
    default:
        currentLevel = LOG_INFO
        log.Printf("[WARN] Invalid LOG_LEVEL '%s'; defaulting to LOG_INFO", levelStr)
    }

    if logFormat == "plain" {
        log.SetFlags(0)
        jsonMode = false
    } else {
        jsonMode = true
        zl = zerolog.New(os.Stdout).With().Timestamp().Logger()
    }
}

func WithRequestID(ctx context.Context, reqID string) context.Context {
    return context.WithValue(ctx, requestIDKey, reqID)
}

func GetRequestID(ctx context.Context) string {
    if ctx == nil {
        return ""
    }
    if reqID, ok := ctx.Value(requestIDKey).(string); ok && reqID != "" {
        return reqID
    }
    return ""
}

// logf is for free-form messages
func logf(ctx context.Context, level int, prefix string, format string, args ...interface{}) {
    if level > currentLevel {
        return
    }

    reqID := GetRequestID(ctx)
    msg := fmt.Sprintf(format, args...)

    if jsonMode {
        e := zl.With().Str("req", reqID).Logger()
        switch level {
        case LOG_ERROR:
            e.Error().Msg(msg)
        case LOG_WARN:
            e.Warn().Msg(msg)
        case LOG_INFO:
            e.Info().Msg(msg)
        case LOG_DEBUG:
            e.Debug().Msg(msg)
        }
    } else {
        timestamp := time.Now().Format("2006/01/02 15:04:05")
        log.Printf("[%s] [req:%s] [%s] %s\n", timestamp, reqID, prefix, msg)
    }
}

func Error(ctx context.Context, format string, args ...interface{}) { logf(ctx, LOG_ERROR, "ERROR", format, args...) }
func Warn(ctx context.Context, format string, args ...interface{})  { logf(ctx, LOG_WARN, "WARN", format, args...) }
func Info(ctx context.Context, format string, args ...interface{})  { logf(ctx, LOG_INFO, "INFO", format, args...) }
func Debug(ctx context.Context, format string, args ...interface{}) { logf(ctx, LOG_DEBUG, "DEBUG", format, args...) }
func Fatal(ctx context.Context, format string, args ...interface{}) { log.Fatalf(format, args...) }

// StructuredRequestLog logs HTTP request details with dedicated fields
func StructuredRequestLog(ctx context.Context, method, path, clientIP string, status int, latency time.Duration) {
    if currentLevel < LOG_INFO {
        return
    }

    reqID := GetRequestID(ctx)

    if jsonMode {
        zl.Info().
            Str("req", reqID).
            Str("method", method).
            Str("path", path).
            Str("client_ip", clientIP).
            Int("status", status).
            Float64("latency_ms", float64(latency.Milliseconds())).
            Msg("http_request")
    } else {
        timestamp := time.Now().Format("2006/01/02 15:04:05")
        log.Printf("[%s] [req:%s] [INFO] %s %s %s %d (%s)", timestamp, reqID, method, path, clientIP, status, latency)
    }
}

func StructuredErrorLog(ctx context.Context, method, path, clientIP string, status int, err error) {
    if currentLevel < LOG_ERROR {
        return
    }

    reqID := GetRequestID(ctx)

    if jsonMode {
        zl.Error().
            Str("req", reqID).
            Str("method", method).
            Str("path", path).
            Str("client_ip", clientIP).
            Int("status", status).
            Str("error", err.Error()).
            Msg("http_error")
    } else {
        timestamp := time.Now().Format("2006/01/02 15:04:05")
        log.Printf("[%s] [req:%s] [ERROR] %s %s %s %d: %v", timestamp, reqID, method, path, clientIP, status, err)
    }
}

/*func StructuredPanicLog(ctx context.Context, method, path string, status int, recovered interface{}, stack []byte) {
    if currentLevel < LOG_ERROR {
        return
    }

    reqID := GetRequestID(ctx)

    if jsonMode {
        zl.Fatal().
            Str("req", reqID).
            Str("method", method).
            Str("path", path).
            Int("status", status).
            Str("panic", fmt.Sprintf("%v", recovered)).
            Str("stack", string(stack)).
            Msg("panic_recovered")
    } else {
        timestamp := time.Now().Format("2006/01/02 15:04:05")
        log.Printf("[%s] [req:%s] [PANIC] %s %s %d - %v\n%s", timestamp, reqID, method, path, status, recovered, string(stack))
    }
}*/
