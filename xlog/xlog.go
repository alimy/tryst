// Copyright 2024 Michael Li <alimy@niubiu.com>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package xlog

import (
	"context"
	"log/slog"
)

var (
	defaultLogger *slog.Logger

	// With calls Logger.With on the default logger.
	With func(args ...any) *slog.Logger

	// WithGroup Logger.WithGroup on the default logger.
	WithGroup func(name string) *slog.Logger

	// Debug calls Logger.Debug on the default logger.
	Debug func(msg string, args ...any)

	// DebugContext calls Logger.DebugContext on the default logger.
	DebugContext func(ctx context.Context, msg string, args ...any)

	// Info calls Logger.Info on the default logger.
	Info func(msg string, args ...any)

	// InfoContext calls Logger.InfoContext on the default logger.
	InfoContext func(ctx context.Context, msg string, args ...any)

	// Warn calls Logger.Warn on the default logger.
	Warn func(msg string, args ...any)

	// WarnContext calls Logger.WarnContext on the default logger.
	WarnContext func(ctx context.Context, msg string, args ...any)

	// Error calls Logger.Error on the default logger.
	Error func(msg string, args ...any)

	// ErrorContext calls Logger.ErrorContext on the default logger.
	ErrorContext func(ctx context.Context, msg string, args ...any)

	// Log calls Logger.Log on the default logger.
	Log func(ctx context.Context, level slog.Level, msg string, args ...any)

	// LogAttrs calls Logger.LogAttrs on the default logger.
	LogAttrs func(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr)
)

func init() {
	SetLogger(slog.Default())
}

// MyLogger returns the default Logger.
func MyLogger() *slog.Logger {
	return defaultLogger
}

// SetLogger makes l the default Logger.
func SetLogger(l *slog.Logger) {
	if l == nil {
		return
	}
	defaultLogger = l

	With, WithGroup = l.With, l.WithGroup
	Debug, DebugContext = l.Debug, l.DebugContext
	Info, InfoContext = l.Info, l.InfoContext
	Warn, WarnContext = l.Warn, l.WarnContext
	Error, ErrorContext = l.Error, l.ErrorContext
	Log, LogAttrs = l.Log, l.LogAttrs
}
