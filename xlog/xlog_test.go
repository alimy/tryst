package xlog_test

import (
	"bytes"
	"context"
	"log/slog"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/alimy/tryst/xlog"
)

const timeRE = `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}(Z|[+-]\d{2}:\d{2})`

func TestLogTextHandler(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer

	l := slog.New(slog.NewTextHandler(&buf, nil))
	xlog.SetLogger(l)

	check := func(want string) {
		t.Helper()
		if want != "" {
			want = "time=" + timeRE + " " + want
		}
		checkLogOutput(t, buf.String(), want)
		buf.Reset()
	}

	l.Info("msg", "a", 1, "b", 2)
	check(`level=INFO msg=msg a=1 b=2`)

	// By default, debug messages are not printed.
	l.Debug("bg", slog.Int("a", 1), "b", 2)
	check("")

	l.Warn("w", slog.Duration("dur", 3*time.Second))
	check(`level=WARN msg=w dur=3s`)

	l.Error("bad", "a", 1)
	check(`level=ERROR msg=bad a=1`)

	l.Log(ctx, slog.LevelWarn+1, "w", slog.Int("a", 1), slog.String("b", "two"))
	check(`level=WARN\+1 msg=w a=1 b=two`)

	l.LogAttrs(ctx, slog.LevelInfo+1, "a b c", slog.Int("a", 1), slog.String("b", "two"))
	check(`level=INFO\+1 msg="a b c" a=1 b=two`)

	l.Info("info", "a", []slog.Attr{slog.Int("i", 1)})
	check(`level=INFO msg=info a.i=1`)

	l.Info("info", "a", slog.GroupValue(slog.Int("i", 1)))
	check(`level=INFO msg=info a.i=1`)
}

func checkLogOutput(t *testing.T, got, wantRegexp string) {
	t.Helper()
	got = clean(got)
	wantRegexp = "^" + wantRegexp + "$"
	matched, err := regexp.MatchString(wantRegexp, got)
	if err != nil {
		t.Fatal(err)
	}
	if !matched {
		t.Errorf("\ngot  %s\nwant %s", got, wantRegexp)
	}
}

// clean prepares log output for comparison.
func clean(s string) string {
	if len(s) > 0 && s[len(s)-1] == '\n' {
		s = s[:len(s)-1]
	}
	return strings.ReplaceAll(s, "\n", "~")
}
