package utils

import (
	"os"
	"testing"
)

var dbg = os.Getenv("PB_DEBUG")

func Debug() bool {
	return (dbg != "")
}

func Assert(condition bool) {
	if !Debug() || condition {
		return
	}
	panic("assert")
}

func Expect(cond bool, t *testing.T, msg string, args ...any) {
	if !cond {
		t.Fatalf(msg, args...)
	}
}
