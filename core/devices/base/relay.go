///////////////////////////////////////////////////////////////////
//
// Future City Project
//
// Copyright (C) 2020-2021 Sergey Denisov. GPLv3
//
// Written by Sergey Denisov aka LittleBuster (DenisovS21@gmail.com)
//
///////////////////////////////////////////////////////////////////

package base

import (
	"github.com/futcity/controller/core/devices"
)

type Relay struct {
	devices.Device
	status   bool
	state    bool
	updateDb func(name string, val bool)
}

func NewRelay(name string, desc string, defStatus bool, updDb func(name string, val bool)) *Relay {
	var dev = &Relay{}

	dev.updateDb = updDb
	dev.SetName(name)
	dev.SetDescription(desc)
	dev.SetOnline(false)
	dev.SetType("relay")
	dev.status = defStatus

	return dev
}

func (r *Relay) SetStatus(value bool) {
	r.status = value
	if r.updateDb != nil {
		go r.updateDb(r.Name(), value)
	}
}

func (r *Relay) Status() bool {
	return r.status
}

func (r *Relay) SetState(value bool) {
	r.state = value
}

func (r *Relay) State() bool {
	return r.state
}

func (r *Relay) Switch() {
	r.SetStatus(!r.Status())
}

func (r *Relay) Update(state bool) {
	r.SetState(state)
	r.SetOnline(true)
}
