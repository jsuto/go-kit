package logx

import (
	"context"
	"testing"
)

func TestWithRequestID(t *testing.T) {
	ctx := WithRequestID(context.Background(), "req-123")
	if got := GetRequestID(ctx); got != "req-123" {
		t.Errorf("GetRequestID = %q, want %q", got, "req-123")
	}
}

func TestGetRequestID_nilContext(t *testing.T) {
	if got := GetRequestID(nil); got != "" {
		t.Errorf("GetRequestID(nil) = %q, want empty string", got)
	}
}

func TestGetRequestID_missing(t *testing.T) {
	if got := GetRequestID(context.Background()); got != "" {
		t.Errorf("GetRequestID(empty ctx) = %q, want empty string", got)
	}
}

func TestWithComponent(t *testing.T) {
	ctx := WithComponent(context.Background(), "purge")
	if got := GetComponent(ctx); got != "purge" {
		t.Errorf("GetComponent = %q, want %q", got, "purge")
	}
}

func TestGetComponent_nilContext(t *testing.T) {
	if got := GetComponent(nil); got != "" {
		t.Errorf("GetComponent(nil) = %q, want empty string", got)
	}
}

func TestGetComponent_missing(t *testing.T) {
	if got := GetComponent(context.Background()); got != "" {
		t.Errorf("GetComponent(empty ctx) = %q, want empty string", got)
	}
}

func TestWithComponentAndRequestID_independent(t *testing.T) {
	ctx := WithRequestID(context.Background(), "req-abc")
	ctx = WithComponent(ctx, "scanner")

	if got := GetRequestID(ctx); got != "req-abc" {
		t.Errorf("GetRequestID = %q, want %q", got, "req-abc")
	}
	if got := GetComponent(ctx); got != "scanner" {
		t.Errorf("GetComponent = %q, want %q", got, "scanner")
	}
}
