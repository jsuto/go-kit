package utils

import (
    "os"
    "regexp"
)

func GetEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

func IsValidParam(param string, pattern string) bool {
    var re = regexp.MustCompile(pattern)
    return re.MatchString(param)
}

func EscapeString(value string) string {
    // Regex to match non-alphanumeric characters
    re := regexp.MustCompile(`[^a-zA-Z0-9_ ]`)
    return re.ReplaceAllStringFunc(value, func(s string) string {
        return "\\" + s
    })
}
