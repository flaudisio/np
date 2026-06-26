package log

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureStderr(t *testing.T, fn func()) string {
	old := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stderr = w
	defer func() { os.Stderr = old }()

	fn()

	_ = w.Close()
	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

func TestDebugEnabled(t *testing.T) {
	DebugEnabled = true
	defer func() { DebugEnabled = false }()

	out := captureStderr(t, func() {
		Debug("debug message")
	})

	if !strings.Contains(out, "DEBUG:") {
		t.Errorf("expected DEBUG prefix, got: %s", out)
	}
	if !strings.Contains(out, "debug message") {
		t.Errorf("expected debug message, got: %s", out)
	}
}

func TestDebugDisabled(t *testing.T) {
	DebugEnabled = false

	out := captureStderr(t, func() {
		Debug("debug message")
	})

	if out != "" {
		t.Errorf("expected no output when debug disabled, got: %s", out)
	}
}

func TestInfo(t *testing.T) {
	out := captureStderr(t, func() {
		Info("hello")
	})

	if out == "" {
		t.Fatal("expected output, got empty")
	}
	if !strings.Contains(out, "INFO:") {
		t.Errorf("expected INFO prefix, got: %s", out)
	}
	if !strings.Contains(out, "hello") {
		t.Errorf("expected hello message, got: %s", out)
	}
}

func TestWarn(t *testing.T) {
	out := captureStderr(t, func() {
		Warn("warning message")
	})

	if !strings.Contains(out, "WARN:") {
		t.Errorf("expected WARN prefix, got: %s", out)
	}
	if !strings.Contains(out, "warning message") {
		t.Errorf("expected warning message, got: %s", out)
	}
}

func TestError(t *testing.T) {
	out := captureStderr(t, func() {
		Error("error message")
	})

	if !strings.Contains(out, "ERROR:") {
		t.Errorf("expected ERROR prefix, got: %s", out)
	}
	if !strings.Contains(out, "error message") {
		t.Errorf("expected error message, got: %s", out)
	}
}

func TestSuccess(t *testing.T) {
	out := captureStderr(t, func() {
		Success("success message")
	})

	if !strings.Contains(out, "SUCCESS:") {
		t.Errorf("expected SUCCESS prefix, got: %s", out)
	}
	if !strings.Contains(out, "success message") {
		t.Errorf("expected success message, got: %s", out)
	}
}

func TestAllLevelsWriteDistinctPrefixes(t *testing.T) {
	levels := map[string]func(string){
		"INFO:":    Info,
		"WARN:":    Warn,
		"ERROR:":   Error,
		"SUCCESS:": Success,
	}

	for prefix, fn := range levels {
		out := captureStderr(t, func() {
			fn("test")
		})

		if !strings.Contains(out, prefix) {
			t.Errorf("expected %s prefix, got: %s", prefix, out)
		}
	}
}
