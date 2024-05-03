package main

import "time"

const (
	TagRegister = "reg"
	TagBeat     = "beat"
	TagPrune    = "prune"
)

type Action struct {
	tag    string
	action any
}

type Register struct {
	id    string
	reply chan<- bool
}

type Beat struct {
	id    string
	ts    time.Time
	reply chan<- bool
}

type Prune struct {
	ts    time.Time
	reply chan<- []string
}

func NewRegister(id string, reply chan<- bool) Action {
	return Action{TagRegister, Register{id, reply}}
}

func NewBeat(id string, t time.Time, reply chan<- bool) Action {
	return Action{TagBeat, Beat{id, t, reply}}
}

func NewPrune(reply chan<- []string) Action {
	return Action{TagPrune, Prune{time.Now(), reply}}
}
