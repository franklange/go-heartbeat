package main

import (
	"testing"
	"time"

	"github.com/franklange/go-heartbeat/utils"
)

func TestCoreAddClient(t *testing.T) {
	core := NewCore()

	added := core.add("c1")
	utils.Expect(added, t, "")

	numClients := len(core.clients)
	utils.Expect(numClients == 1, t, "numClients: %d", numClients)

	numTs := len(core.clients["c1"])
	utils.Expect(numTs == 0, t, "numTs: %d", numTs)
}

func TestCoreAddClientRedundant(t *testing.T) {
	core := NewCore()

	added := core.add("c1")
	utils.Expect(added, t, "")

	added = core.add("c1")
	utils.Expect(!added, t, "")

	numClients := len(core.clients)
	utils.Expect(numClients == 1, t, "numClients: %d", numClients)

	numTs := len(core.clients["c1"])
	utils.Expect(numTs == 0, t, "numTs: %d", numTs)
}

func TestCoreBeatUnknownClient(t *testing.T) {
	core := NewCore()

	added := core.beat("c1")
	utils.Expect(!added, t, "")
}

func TestCoreBeatFirst(t *testing.T) {
	core := NewCore()
	core.add("c1")

	added := core.beat("c1")
	utils.Expect(added, t, "")

	numTs := len(core.clients["c1"])
	utils.Expect(numTs == 1, t, "numTs: %d", numTs)
}

func TestCoreBeatSingleClientMuliBeats(t *testing.T) {
	core := NewCore()

	core.add("c1")
	core.beat_at("c1", time.Now())
	core.beat_at("c1", time.Now().Add(1*time.Second))
	core.beat_at("c1", time.Now().Add(2*time.Second))

	numTs := len(core.clients["c1"])
	utils.Expect(numTs == 3, t, "numTs: %d", numTs)
}

func TestCoreBeatMultiClientMuliBeats(t *testing.T) {
	core := NewCore()

	core.add("c1")
	core.add("c2")

	core.beat_at("c1", time.Now())
	core.beat_at("c1", time.Now().Add(1*time.Second))
	core.beat_at("c1", time.Now().Add(2*time.Second))

	core.beat_at("c2", time.Now())
	core.beat_at("c2", time.Now().Add(2*time.Second))

	numClients := len(core.clients)
	utils.Expect(numClients == 2, t, "numClients: %d", numClients)

	numTs1 := len(core.clients["c1"])
	utils.Expect(numTs1 == 3, t, "numTs: %d", numTs1)

	numTs2 := len(core.clients["c2"])
	utils.Expect(numTs2 == 2, t, "numTs: %d", numTs2)
}

func TestCoreBeatOutdated(t *testing.T) {
	core := NewCore()
	core.add("c1")
	core.beat_at("c1", time.Now().Add(5*time.Second))
	core.beat_at("c1", time.Now().Add(6*time.Second))

	added := core.beat_at("c1", time.Now())
	utils.Expect(!added, t, "add outdated")

	numTs := len(core.clients["c1"])
	utils.Expect(numTs == 2, t, "numTs: %d", numTs)
}

func TestCorePruneEmpty(t *testing.T) {
	core := NewCore()
	res := core.prune()

	numDead := len(res)
	utils.Expect(numDead == 0, t, "numDead: %d", numDead)
}

func TestCorePruneSingleClientNoTs(t *testing.T) {
	core := NewCore()
	core.add("c1")

	res := core.prune()

	numDead := len(res)
	utils.Expect(numDead == 1, t, "numDead: %d", numDead)

	numClients := len(core.clients)
	utils.Expect(numClients == 0, t, "numClients: %d", numClients)
}

func TestCorePruneMultiClientNoTs(t *testing.T) {
	core := NewCore()
	core.add("c1")
	core.add("c2")

	res := core.prune()

	numDead := len(res)
	utils.Expect(numDead == 2, t, "numDead: %d", numDead)

	numClients := len(core.clients)
	utils.Expect(numClients == 0, t, "numClients: %d", numClients)
}

func TestCorePruneInThePast(t *testing.T) {
	core := NewCore()
	core.add("c1")

	core.beat_at("c1", time.Now().Add(5*time.Second))
	core.beat_at("c1", time.Now().Add(6*time.Second))

	res := core.prune()

	numDead := len(res)
	utils.Expect(numDead == 0, t, "numDead: %d", numDead)

	numClients := len(core.clients)
	utils.Expect(numClients == 1, t, "numClients: %d", numClients)
}

func TestCorePruneSingleClient(t *testing.T) {
	core := NewCore()
	core.add("c1")
	core.beat_at("c1", time.Now())
	core.beat_at("c1", time.Now().Add(1*time.Second))

	res := core.prune_at(time.Now().Add(10 * time.Second))

	numDead := len(res)
	utils.Expect(numDead == 1, t, "numDead: %d", numDead)

	numClients := len(core.clients)
	utils.Expect(numClients == 0, t, "numClients: %d", numClients)
}

func TestCorePruneMultiClientAllDead(t *testing.T) {
	core := NewCore()
	core.add("c1")
	core.beat_at("c1", time.Now())
	core.beat_at("c1", time.Now().Add(1*time.Second))

	core.add("c2")
	core.beat_at("c2", time.Now())
	core.beat_at("c2", time.Now().Add(1*time.Second))

	res := core.prune_at(time.Now().Add(10 * time.Second))

	numDead := len(res)
	utils.Expect(numDead == 2, t, "numDead: %d", numDead)

	numClients := len(core.clients)
	utils.Expect(numClients == 0, t, "numClients: %d", numClients)
}
