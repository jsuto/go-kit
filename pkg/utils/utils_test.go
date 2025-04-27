package utils

import (
    "testing"
)

func TestEscapeString(t *testing.T) {
    tests := []struct {
        input    string
        expected string
    }{
        {"simpleTest123", "simpleTest123"},         // Alphanumeric - no change
        {"test'quote", "test\\'quote"},             // Single quote
        {"hello;drop", "hello\\;drop"},             // Semicolon
        {"newline\n", "newline\\\n"},               // Newline character
        {"special#chars!", "special\\#chars\\!"},   // Special characters
        {"mixed'\"chars;", "mixed\\'\\\"chars\\;"}, // Mixed characters
    }

    for _, tt := range tests {
        t.Run(tt.input, func(t *testing.T) {
            result := EscapeString(tt.input)
            if result != tt.expected {
                t.Errorf("escapeString(%q) = %q; want %q", tt.input, result, tt.expected)
            }
        })
    }
}

func TestIsValidParam(t *testing.T) {
    tests := []struct {
        param     string
        pattern   string
        expected  bool
        testName  string
    }{
        {"validParam123", `^[a-zA-Z0-9]+$`, true, "AlphanumericOnly"},     // Valid alphanumeric
        {"invalid@char", `^[a-zA-Z0-9]+$`, false, "InvalidCharacter"},      // Invalid character
        {"123", `^\d+$`, true, "OnlyNumbers"},                              // Numeric only
        {"123abc", `^\d+$`, false, "AlphanumericWithNumbers"},              // Pattern only allows numbers
        {"UPPERCASE", `^[A-Z]+$`, true, "UpperCaseOnly"},                   // Uppercase only
        {"lowercase", `^[A-Z]+$`, false, "UppercaseOnlyWithLowercaseInput"}, // Pattern allows uppercase only
    }

    for _, tt := range tests {
        t.Run(tt.testName, func(t *testing.T) {
            result := IsValidParam(tt.param, tt.pattern)
            if result != tt.expected {
                t.Errorf("isValidParam(%q, %q) = %v; want %v", tt.param, tt.pattern, result, tt.expected)
            }
        })
    }
}
