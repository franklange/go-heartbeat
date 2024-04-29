package main

const (
	TagRegister  = "reg"
	TagHeartbeat = "beat"
	TagPrune     = "prune"
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
	reply chan<- bool
}

type Prune struct {
	reply chan<- []string
}

type Update struct {
	tag    string
	update any
}

type RegisterUpdate struct {
	id string
}

type PruneUpdate struct {
	ids []string
}
