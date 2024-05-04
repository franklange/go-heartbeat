package heartbeat

import (
	"slices"
	"strings"
	"testing"
	"time"
)

func expect(cond bool, t *testing.T, msg string, args ...any) {
	if !cond {
		t.Fatalf(msg, args...)
	}
}

func advanceS(n int) time.Time {
	return time.Now().Add(time.Duration(n) * time.Second)
}

func split(ids string) []string {
	if len(ids) == 0 {
		return []string{}
	}
	return strings.Split(strings.ReplaceAll(ids, " ", ""), ",")
}

func (h *Heartbeats) hasClients(ids string) bool {
	clients := split(ids)
	if len(clients) != h.NumClients() {
		return false
	}
	for _, v := range clients {
		if !h.HasClient(v) {
			return false
		}
	}
	return true
}

func eq(a, b []string) bool {
	slices.Sort(a)
	slices.Sort(b)
	return (slices.Compare(a, b) == 0)
}

func newHeartbeats() *Heartbeats {
	return NewHeartbeats(5 * time.Second)
}

func TestRegisterSingleClient(t *testing.T) {
	h := newHeartbeats()

	ok := h.Register("c1")
	expect(ok && h.hasClients("c1"), t, "register")
}

func TestRegisterSameClient(t *testing.T) {
	h := newHeartbeats()

	ok := h.Register("c1")
	expect(ok && h.hasClients("c1"), t, "register")

	ok = h.Register("c1")
	expect(!ok && h.hasClients("c1"), t, "re-register")
}

func TestRegisterMultiClient(t *testing.T) {
	h := newHeartbeats()

	ok := h.Register("c1")
	expect(ok && h.hasClients("c1"), t, "first client")

	ok = h.Register("c2")
	expect(ok && h.hasClients("c1, c2"), t, "second client")
}

func TestBeatUnknownClient(t *testing.T) {
	h := newHeartbeats()

	ok := h.Beat("c1")
	expect(!ok && h.hasClients(""), t, "unknown client")
}

func TestBeatSingleClient(t *testing.T) {
	h := newHeartbeats()
	h.Register("c1")

	ok := h.beat_at("c1", advanceS(1))
	expect(ok && h.hasClients("c1"), t, "beat")
}

func TestBeatSingleClientMultiBeats(t *testing.T) {
	h := newHeartbeats()
	h.Register("c1")

	ok := h.beat_at("c1", advanceS(1))
	expect(ok && h.hasClients("c1"), t, "first beat")

	ok = h.beat_at("c1", advanceS(1))
	expect(ok && h.hasClients("c1"), t, "second beat")
}

func TestBeatMulitClientMultiBeats(t *testing.T) {
	h := newHeartbeats()
	h.Register("c1")
	h.Register("c2")

	ok := h.beat_at("c1", advanceS(1))
	expect(ok && h.hasClients("c1, c2"), t, "first c1")

	ok = h.beat_at("c2", advanceS(2))
	expect(ok && h.hasClients("c1, c2"), t, "first c2")

	ok = h.beat_at("c1", advanceS(3))
	expect(ok && h.hasClients("c1, c2"), t, "second c1")
}

func TestBeatOutdated(t *testing.T) {
	h := newHeartbeats()
	h.register_at("c1", advanceS(3))

	ok := h.Beat("c1")
	expect(!ok && h.hasClients("c1"), t, "outdated beat")
}

func TestPruneEmpty(t *testing.T) {
	h := newHeartbeats()

	have := h.Prune()
	want := []string{}

	expect(eq(have, want), t, "prune empty")
}

func TestPruneSingleClientDead(t *testing.T) {
	h := newHeartbeats()

	h.Register("c1")
	h.beat_at("c1", advanceS(1))

	have := h.prune_at(advanceS(10))
	want := []string{"c1"}

	expect(eq(have, want), t, "single client dead")
}

func TestPruneSingleClientAlive(t *testing.T) {
	h := newHeartbeats()

	h.Register("c1")
	h.beat_at("c1", advanceS(1))

	have := h.prune_at(advanceS(2))
	want := []string{}

	expect(eq(have, want), t, "single client alive")
}

func TestPruneMultiClientAllDead(t *testing.T) {
	h := newHeartbeats()

	h.Register("c1")
	h.Register("c2")

	have := h.prune_at(advanceS(12))
	want := []string{"c1", "c2"}

	expect(eq(have, want), t, "prune all dead")
}

func TestPruneMultiClientSomeDead(t *testing.T) {
	h := newHeartbeats()

	h.Register("c1")
	h.Register("c2")
	h.beat_at("c1", advanceS(3))

	have := h.prune_at(advanceS(6))
	want := []string{"c2"}

	expect(eq(have, want), t, "prune some dead")
}

func TestPruneInThePast(t *testing.T) {
	h := newHeartbeats()

	h.register_at("c1", advanceS(10))
	h.register_at("c2", advanceS(11))

	have := h.Prune()
	want := []string{}

	expect(eq(have, want), t, "prune in the past")
}
