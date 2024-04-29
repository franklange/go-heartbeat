package main

import (
	"fmt"
	"log/slog"
)

type Service struct {
	core    Core
	actions chan Action
	updates chan Update
}

func NewService() Service {
	return Service{NewCore(), make(chan Action, 10), make(chan Update, 10)}
}

func (service *Service) runOne() {
	a := <-service.actions
	slog.Debug(fmt.Sprint(a))
	service.route(a)
}

func (service *Service) runAll() {
	for len(service.actions) > 0 {
		service.runOne()
	}
}

func (service *Service) runForever() {
	for {
		service.runOne()
	}
}

func (service *Service) route(a Action) {
	if a.tag == TagRegister {
		service.register(a.action.(Register))
		return
	}

	if a.tag == TagHeartbeat {
		service.heartbeat(a.action.(Beat))
		return
	}

	if a.tag == TagPrune {
		service.prune(a.action.(Prune))
		return
	}
	slog.Debug("unknown action")
}

func (service *Service) register(a Register) {
	ok := service.core.add(a.id)
	if !ok {
		slog.Debug("register failed")
	} else {
		service.updates <- Update{TagRegister, RegisterUpdate{a.id}}
	}
	a.reply <- ok
}

func (service *Service) heartbeat(a Beat) {
	ok := service.core.beat(a.id)
	if !ok {
		slog.Debug("heartbeat failed")
	}
	a.reply <- ok
}

func (service *Service) prune(a Prune) {
	a.reply <- service.core.prune()
}
