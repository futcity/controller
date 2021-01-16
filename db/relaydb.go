///////////////////////////////////////////////////////////////////
//
// Future City Project
//
// Copyright (C) 2020-2021 Sergey Denisov. GPLv3
//
// Written by Sergey Denisov aka LittleBuster (DenisovS21@gmail.com)
//
///////////////////////////////////////////////////////////////////

package db

type SingleRelayDB struct {
	Name   string `json:"name"`
	Status bool   `json:"status"`
}

type RelaysDB struct {
	Relays []SingleRelayDB `json:"relays"`
}
