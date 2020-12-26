///////////////////////////////////////////////////////////////////
//
// Future City Project
//
// Copyright (C) 2020-2021 Sergey Denisov. GPLv3
//
// Written by Sergey Denisov aka LittleBuster (DenisovS21@gmail.com)
//
///////////////////////////////////////////////////////////////////

package auth

type ProfileDevice struct {
	name  string
	read  bool
	write bool
}

func NewProfileDevice(name string, read bool, write bool) *ProfileDevice {
	return &ProfileDevice{
		name:  name,
		read:  read,
		write: write,
	}
}

func (p *ProfileDevice) Name() string {
	return p.name
}

func (p *ProfileDevice) Read() bool {
	return p.read
}

func (p *ProfileDevice) Write() bool {
	return p.write
}
