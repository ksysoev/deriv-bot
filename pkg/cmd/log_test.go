package cmd

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockHandler is a mock implementation of slog.Handler for testing
type MockHandler struct {
	capturedAttrs map[string]string
}

func NewMockHandler() *MockHandler {
	return &MockHandler{
		capturedAttrs: make(map[string]string),
	}
}

func (h *MockHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

//nolint:gocritic // this is a mock handler, so we don't need to implement all methods
func (h *MockHandler) Handle(_ context.Context, r slog.Record) error {
	r.Attrs(func(attr slog.Attr) bool {
		if attr.Value.Kind() == slog.KindString {
			h.capturedAttrs[attr.Key] = attr.Value.String()
		}

		return true
	})

	return nil
}

func (h *MockHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h *MockHandler) WithGroup(_ string) slog.Handler {
	return h
}

func TestContextHandler_Handle(t *testing.T) {
	//nolint:staticcheck // don't want to have dependecy on cmd package here for now
	tests := []struct {
		ctx       context.Context
		wantAttrs map[string]string
		name      string
		app       string
		version   string
	}{
		{
			name:    "With_request_ID",
			ctx:     context.WithValue(context.Background(), "req_id", "test-request-id"),
			app:     "test-app",
			version: "1.0.0",
			wantAttrs: map[string]string{
				"req_id": "test-request-id",
				"app":    "test-app",
				"ver":    "1.0.0",
			},
		},
		{
			name:    "Without_request_ID",
			ctx:     context.Background(),
			app:     "test-app",
			version: "1.0.0",
			wantAttrs: map[string]string{
				"app": "test-app",
				"ver": "1.0.0",
			},
		},
		{
			name:    "Handles_no_attributes_case",
			ctx:     context.Background(),
			app:     "",
			version: "",
			wantAttrs: map[string]string{
				"app": "",
				"ver": "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHandler := NewMockHandler()
			ch := ContextHandler{
				Handler: mockHandler,
				app:     tt.app,
				ver:     tt.version,
			}

			// Create a simple log record
			record := slog.Record{
				Time:    time.Now(),
				Message: "test message",
				Level:   slog.LevelInfo,
			}

			err := ch.Handle(tt.ctx, record)
			assert.NoError(t, err)

			// Check if all expected attributes are present with correct values
			for k, v := range tt.wantAttrs {
				assert.Equal(t, v, mockHandler.capturedAttrs[k])
			}
		})
	}
}

func TestInitLogger(t *testing.T) {
	tests := []struct {
		name    string
		arg     cmdArgs
		wantErr bool
	}{
		{
			name: "Valid log level with text format",
			arg: cmdArgs{
				LogLevel:   "info",
				TextFormat: true,
				Version:    "1.0.0",
			},
			wantErr: false,
		},
		{
			name: "Valid log level with JSON format",
			arg: cmdArgs{
				LogLevel:   "debug",
				TextFormat: false,
				Version:    "1.0.0",
			},
			wantErr: false,
		},
		{
			name: "Invalid log level",
			arg: cmdArgs{
				LogLevel:   "invalid",
				TextFormat: true,
				Version:    "1.0.0",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.arg
			err := initLogger(&args)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
