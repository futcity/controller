///////////////////////////////////////////////////////////////////
//
// Future City Project
//
// Copyright (C) 2020-2021 Sergey Denisov. GPLv3
//
// Written by Sergey Denisov aka LittleBuster (DenisovS21@gmail.com)
//
///////////////////////////////////////////////////////////////////

package devices

type IDevice interface {
	ID() int
	SetID(id int)
	Name() string
	SetName(name string)
	SetType(value string)
	Type() string
	SetDescription(key string)
	Description() string
	SetOnline(value bool)
	Online() bool
}

type Device struct {
	id      int
	name    string
	online  bool
	desc    string
	devType string
}

func (d *Device) ID() int {
	return d.id
}

func (d *Device) SetID(id int) {
	d.id = id
}

func (d *Device) Name() string {
	return d.name
}

func (d *Device) SetName(name string) {
	d.name = name
}

func (d *Device) Online() bool {
	return d.online
}

func (d *Device) SetOnline(value bool) {
	d.online = value
}

func (d *Device) Description() string {
	return d.desc
}

func (d *Device) SetDescription(desc string) {
	d.desc = desc
}

func (d *Device) Type() string {
	return d.devType
}

func (d *Device) SetType(value string) {
	d.devType = value
}
