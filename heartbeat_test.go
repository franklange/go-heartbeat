package heartbeat

import (
	"slices"
	"testing"
)

func expect(cond bool, t *testing.T, msg string, args ...any) {
	if !cond {
		t.Fatalf(msg, args...)
	}
}

func eq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	slices.Sort(a)
	slices.Sort(b)
	return (slices.Compare(a, b) == 0)
}

func TestHasClientEmpty(t *testing.T) {
	h := NewHeartbeats()

	have := h.HasClient("c1")
	want := false

	expect(have == want, t, "has client empty")
}

func TestHasClient(t *testing.T) {
	h := NewHeartbeats()
	h.Beat("c1")

	have := h.HasClient("c1")
	want := true

	expect(have == want, t, "has client")
}

func TestBeatEmpty(t *testing.T) {
	h := NewHeartbeats()
	h.Beat("c1")

	have := h.numBeats("c1")
	want := 3

	expect(have == want, t, "beat init")
}

func TestBeatIncrements(t *testing.T) {
	h := NewHeartbeats()
	h.clients["c1"] = 1
	h.Beat("c1")

	have := h.numBeats("c1")
	want := 2

	expect(have == want, t, "beat increment")
}

func TestBeatCap(t *testing.T) {
	h := NewHeartbeats()
	h.Beat("c1")
	h.Beat("c1")

	have := h.numBeats("c1")
	want := 3

	expect(have == want, t, "beat cap")
}

func TestBeatMultiClient(t *testing.T) {
	h := NewHeartbeats()
	h.clients["c1"] = 1
	h.clients["c2"] = 2
	h.Beat("c1")

	have := h.numBeats("c1")
	want := 2
	expect(have == want, t, "beat multi client inc")

	have = h.numBeats("c2")
	want = 2
	expect(have == want, t, "beat multi client no change")
}

func TestPruneEmpty(t *testing.T) {
	h := NewHeartbeats()

	have := h.Prune()
	want := []string{}

	expect(eq(have, want), t, "prune empty")
}

func TestPruneSingleClientAlive(t *testing.T) {
	h := NewHeartbeats()
	h.Beat("c1")

	have := h.Prune()
	want := []string{}

	expect(eq(have, want), t, "prune single alive")
}

func TestPruneSingleClientDead(t *testing.T) {
	h := NewHeartbeats()
	h.clients["c1"] = 2
	h.Prune()

	have := h.Prune()
	want := []string{"c1"}

	expect(eq(have, want), t, "prune single dead")

}

func TestPruneMultiClientAllAlice(t *testing.T) {
	h := NewHeartbeats()
	h.Beat("c1")
	h.Beat("c2")

	have := h.Prune()
	want := []string{}

	expect(eq(have, want), t, "prune multi all alive")
}

func TestPruneMultiClientAllDead(t *testing.T) {
	h := NewHeartbeats()
	h.clients["c1"] = 0
	h.clients["c2"] = 1

	have := h.Prune()
	want := []string{"c1", "c2"}

	expect(eq(have, want), t, "prune multi all dead")
}

func TestPruneMultiClientSomeDead(t *testing.T) {
	h := NewHeartbeats()
	h.clients["c1"] = 2
	h.clients["c2"] = 1

	have := h.Prune()
	want := []string{"c2"}

	expect(eq(have, want), t, "prune multi some dead")
}
