package lib

import (
	"testing"
	"time"
)

type CoreTest struct {
	core       *Core
	regreply   chan bool
	beatreply  chan bool
	prunereply chan []string
}

func NewCoreTest() CoreTest {
	return CoreTest{NewCore(), make(chan bool, 1), make(chan bool, 1), make(chan []string, 1)}
}

func (ct *CoreTest) register(id string) bool {
	ct.core.actions <- newRegister(id, ct.regreply)
	ct.core.runOne()
	return <-ct.regreply
}

func (ct *CoreTest) beat(id string, t time.Time) bool {
	ct.core.actions <- newBeat(id, t, ct.beatreply)
	ct.core.runOne()
	return <-ct.beatreply
}

func (ct *CoreTest) prune(t time.Time) []string {
	ct.core.actions <- newPrune(t, ct.prunereply)
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
	expect(ok, t, "")

	numClients := ct.numClients()
	expect(numClients == 1, t, "numClients: %d", numClients)

	numTs := ct.numTimestamps("c1")
	expect(numTs == 0, t, "numTs: %d", numTs)
}

func TestCoreAddClientRedundant(t *testing.T) {
	ct := NewCoreTest()

	ok := ct.register("c1")
	expect(ok, t, "add first client")

	ok = ct.register("c1")
	expect(!ok, t, "add second client")

	numClients := ct.numClients()
	expect(numClients == 1, t, "numClients: %d", numClients)

	numTs := ct.numTimestamps("c1")
	expect(numTs == 0, t, "numTs: %d", numTs)
}

func TestCoreBeatUnknownClient(t *testing.T) {
	ct := NewCoreTest()

	ok := ct.beat("c1", time.Now())
	expect(!ok, t, "beat unknown client")
}

func TestCoreBeatFirst(t *testing.T) {
	ct := NewCoreTest()

	ct.register("c1")
	ok := ct.beat("c1", time.Now())
	expect(ok, t, "first beat")

	numTs := ct.numTimestamps("c1")
	expect(numTs == 1, t, "numTs: %d", numTs)
}

func TestCoreBeatSingleClientMuliBeats(t *testing.T) {
	ct := NewCoreTest()

	ct.register("c1")
	ct.beat("c1", time.Now())
	ct.beat("c1", time.Now().Add(1*time.Second))
	ct.beat("c1", time.Now().Add(2*time.Second))

	numTs := ct.numTimestamps("c1")
	expect(numTs == 3, t, "numTs: %d", numTs)
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
	expect(numClients == 2, t, "numClients: %d", numClients)

	numTs1 := ct.numTimestamps("c1")
	expect(numTs1 == 3, t, "numTs: %d", numTs1)

	numTs2 := ct.numTimestamps("c2")
	expect(numTs2 == 2, t, "numTs: %d", numTs2)
}

func TestCoreBeatOutdated(t *testing.T) {
	ct := NewCoreTest()
	ct.register("c1")
	ct.beat("c1", time.Now().Add(5*time.Second))
	ct.beat("c1", time.Now().Add(6*time.Second))

	ok := ct.beat("c1", time.Now())
	expect(!ok, t, "add outdated")

	numTs := ct.numTimestamps("c1")
	expect(numTs == 2, t, "numTs: %d", numTs)
}

func TestCorePruneEmpty(t *testing.T) {
	ct := NewCoreTest()
	res := ct.prune(time.Now())

	numDead := len(res)
	expect(numDead == 0, t, "numDead: %d", numDead)
}

func TestCorePruneSingleClientNoTs(t *testing.T) {
	ct := NewCoreTest()
	ct.register("c1")

	res := ct.prune(time.Now())

	numDead := len(res)
	expect(numDead == 1, t, "numDead: %d", numDead)

	numClients := ct.numClients()
	expect(numClients == 0, t, "numClients: %d", numClients)
}

func TestCorePruneMultiClientNoTs(t *testing.T) {
	ct := NewCoreTest()
	ct.register("c1")
	ct.register("c2")

	res := ct.prune(time.Now())

	numDead := len(res)
	expect(numDead == 2, t, "numDead: %d", numDead)

	numClients := ct.numClients()
	expect(numClients == 0, t, "numClients: %d", numClients)
}

func TestCorePruneInThePast(t *testing.T) {
	ct := NewCoreTest()
	ct.register("c1")

	ct.beat("c1", time.Now().Add(5*time.Second))
	ct.beat("c1", time.Now().Add(6*time.Second))

	res := ct.prune(time.Now())

	numDead := len(res)
	expect(numDead == 0, t, "numDead: %d", numDead)

	numClients := ct.numClients()
	expect(numClients == 1, t, "numClients: %d", numClients)
}

func TestCorePruneSingleClient(t *testing.T) {
	ct := NewCoreTest()
	ct.register("c1")
	ct.beat("c1", time.Now())
	ct.beat("c1", time.Now().Add(1*time.Second))

	res := ct.prune(time.Now().Add(10 * time.Second))

	numDead := len(res)
	expect(numDead == 1, t, "numDead: %d", numDead)

	numClients := ct.numClients()
	expect(numClients == 0, t, "numClients: %d", numClients)
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
	expect(numDead == 2, t, "numDead: %d", numDead)

	numClients := ct.numClients()
	expect(numClients == 0, t, "numClients: %d", numClients)
}
