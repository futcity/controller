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

import "github.com/futcity/controller/core/devices"

type Light struct {
	devices.Device
	status   bool
	state    bool
	updateDb func(name string, val bool)
}

func NewLight(name string, desc string, defStatus bool, updDb func(name string, val bool)) *Light {
	var dev = &Light{}

	dev.updateDb = updDb
	dev.SetName(name)
	dev.SetDescription(desc)
	dev.SetOnline(false)
	dev.SetType("light")
	dev.status = defStatus

	return dev
}

func (r *Light) SetStatus(value bool) {
	r.status = value
	if r.updateDb != nil {
		go r.updateDb(r.Name(), value)
	}
}

func (r *Light) Status() bool {
	return r.status
}

func (r *Light) SetState(value bool) {
	r.state = value
}

func (r *Light) State() bool {
	return r.state
}

func (r *Light) Switch() {
	r.SetStatus(!r.Status())
}

func (r *Light) Update(state bool) {
	r.SetState(state)
	r.SetOnline(true)
}
