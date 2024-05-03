package main

import (
	"testing"
	"time"

	"github.com/franklange/go-heartbeat/utils"
)

type CoreTest struct {
	core       Core
	regreply   chan bool
	beatreply  chan bool
	prunereply chan []string
}

func NewCoreTest() CoreTest {
	return CoreTest{NewCore(), make(chan bool, 1), make(chan bool, 1), make(chan []string, 1)}
}

func (ct *CoreTest) NewRegister(id string) Action {
	return Action{TagRegister, Register{id, ct.regreply}}
}

func (ct *CoreTest) NewBeat(id string, t time.Time) Action {
	return Action{TagBeat, Beat{id, t, ct.beatreply}}
}

func (ct *CoreTest) NewPrune(t time.Time) Action {
	return Action{TagPrune, Prune{t, ct.prunereply}}
}

func (ct *CoreTest) register(id string) bool {
	ct.core.actions <- ct.NewRegister(id)
	ct.core.runOne()
	return <-ct.regreply
}

func (ct *CoreTest) beat(id string, t time.Time) bool {
	ct.core.actions <- ct.NewBeat(id, t)
	ct.core.runOne()
	return <-ct.beatreply
}

func (ct *CoreTest) prune(t time.Time) []string {
	ct.core.actions <- ct.NewPrune(t)
	ct.core.runOne()
	return <-ct.prunereply
}

func (ct *CoreTest) numClients() int {
	return len(ct.core.clients)
}

func (ct *CoreTest) numTimestamps(client string) int {
	return len(ct.core.clients[client])
}

func TestCoreAddClient(t *testing.T) {
	ct := NewCoreTest()

	ok := ct.register("c1")
	utils.Expect(ok, t, "")

	numClients := ct.numClients()
	utils.Expect(numClients == 1, t, "numClients: %d", numClients)

	numTs := ct.numTimestamps("c1")
	utils.Expect(numTs == 0, t, "numTs: %d", numTs)
}

func TestCoreAddClientRedundant(t *testing.T) {
	ct := NewCoreTest()

	ok := ct.register("c1")
	utils.Expect(ok, t, "add first client")

	ok = ct.register("c1")
	utils.Expect(!ok, t, "add second client")

	numClients := ct.numClients()
	utils.Expect(numClients == 1, t, "numClients: %d", numClients)

	numTs := ct.numTimestamps("c1")
	utils.Expect(numTs == 0, t, "numTs: %d", numTs)
}

func TestCoreBeatUnknownClient(t *testing.T) {
	ct := NewCoreTest()

	ok := ct.beat("c1", time.Now())
	utils.Expect(!ok, t, "beat unknown client")
}

func TestCoreBeatFirst(t *testing.T) {
	ct := NewCoreTest()

	ct.register("c1")
	ok := ct.beat("c1", time.Now())
	utils.Expect(ok, t, "first beat")

	numTs := ct.numTimestamps("c1")
	utils.Expect(numTs == 1, t, "numTs: %d", numTs)
}

func TestCoreBeatSingleClientMuliBeats(t *testing.T) {
	ct := NewCoreTest()

	ct.register("c1")
	ct.beat("c1", time.Now())
	ct.beat("c1", time.Now().Add(1*time.Second))
	ct.beat("c1", time.Now().Add(2*time.Second))

	numTs := ct.numTimestamps("c1")
	utils.Expect(numTs == 3, t, "numTs: %d", numTs)
}

func TestCoreBeatMultiClientMuliBeats(t *testing.T) {
	ct := NewCoreTest()

	ct.register("c1")
	ct.register("c2")

	ct.beat("c1", time.Now())
	ct.beat("c1", time.Now().Add(1*time.Second))
	ct.beat("c1", time.Now().Add(2*time.Second))

	ct.beat("c2", time.Now())
	ct.beat("c2", time.Now().Add(2*time.Second))

	numClients := ct.numClients()
	utils.Expect(numClients == 2, t, "numClients: %d", numClients)

	numTs1 := ct.numTimestamps("c1")
	utils.Expect(numTs1 == 3, t, "numTs: %d", numTs1)

	numTs2 := ct.numTimestamps("c2")
	utils.Expect(numTs2 == 2, t, "numTs: %d", numTs2)
}

func TestCoreBeatOutdated(t *testing.T) {
	ct := NewCoreTest()
	ct.register("c1")
	ct.beat("c1", time.Now().Add(5*time.Second))
	ct.beat("c1", time.Now().Add(6*time.Second))

	ok := ct.beat("c1", time.Now())
	utils.Expect(!ok, t, "add outdated")

	numTs := ct.numTimestamps("c1")
	utils.Expect(numTs == 2, t, "numTs: %d", numTs)
}

func TestCorePruneEmpty(t *testing.T) {
	ct := NewCoreTest()
	res := ct.prune(time.Now())

	numDead := len(res)
	utils.Expect(numDead == 0, t, "numDead: %d", numDead)
}

func TestCorePruneSingleClientNoTs(t *testing.T) {
	ct := NewCoreTest()
	ct.register("c1")

	res := ct.prune(time.Now())

	numDead := len(res)
	utils.Expect(numDead == 1, t, "numDead: %d", numDead)

	numClients := ct.numClients()
	utils.Expect(numClients == 0, t, "numClients: %d", numClients)
}

func TestCorePruneMultiClientNoTs(t *testing.T) {
	ct := NewCoreTest()
	ct.register("c1")
	ct.register("c2")

	res := ct.prune(time.Now())

	numDead := len(res)
	utils.Expect(numDead == 2, t, "numDead: %d", numDead)

	numClients := ct.numClients()
	utils.Expect(numClients == 0, t, "numClients: %d", numClients)
}

func TestCorePruneInThePast(t *testing.T) {
	ct := NewCoreTest()
	ct.register("c1")

	ct.beat("c1", time.Now().Add(5*time.Second))
	ct.beat("c1", time.Now().Add(6*time.Second))

	res := ct.prune(time.Now())

	numDead := len(res)
	utils.Expect(numDead == 0, t, "numDead: %d", numDead)

	numClients := ct.numClients()
	utils.Expect(numClients == 1, t, "numClients: %d", numClients)
}

func TestCorePruneSingleClient(t *testing.T) {
	ct := NewCoreTest()
	ct.register("c1")
	ct.beat("c1", time.Now())
	ct.beat("c1", time.Now().Add(1*time.Second))

	res := ct.prune(time.Now().Add(10 * time.Second))

	numDead := len(res)
	utils.Expect(numDead == 1, t, "numDead: %d", numDead)

	numClients := ct.numClients()
	utils.Expect(numClients == 0, t, "numClients: %d", numClients)
}

func TestCorePruneMultiClientAllDead(t *testing.T) {
	ct := NewCoreTest()
	ct.register("c1")
	ct.beat("c1", time.Now())
	ct.beat("c1", time.Now().Add(1*time.Second))

	ct.register("c2")
	ct.beat("c2", time.Now())
	ct.beat("c2", time.Now().Add(1*time.Second))

	res := ct.prune(time.Now().Add(10 * time.Second))

	numDead := len(res)
	utils.Expect(numDead == 2, t, "numDead: %d", numDead)

	numClients := ct.numClients()
	utils.Expect(numClients == 0, t, "numClients: %d", numClients)
}
