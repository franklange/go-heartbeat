package lib

import (
	"os"
	"testing"
)

var dbg = os.Getenv("HB_DEBUG")

func Debug() bool {
	return (dbg != "")
}

func assert(condition bool) {
	if !Debug() || condition {
		return
	}
	panic("assert")
}

func expect(cond bool, t *testing.T, msg string, args ...any) {
	if !cond {
		t.Fatalf(msg, args...)
	}
}
