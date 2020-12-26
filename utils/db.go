///////////////////////////////////////////////////////////////////
//
// Future City Project
//
// Copyright (C) 2020-2021 Sergey Denisov. GPLv3
//
// Written by Sergey Denisov aka LittleBuster (DenisovS21@gmail.com)
//
///////////////////////////////////////////////////////////////////

package utils

const (
	DbTextType = iota
	DbRedisType
)

type SingleRelayDb struct {
	Name   string
	Status bool
}

type RelayDb struct {
	Relays []SingleRelayDb
}

type Database struct {
	cfg       *Configs
	fileNames map[string]string
	dbType    int
	relays    RelayDb
}

func NewDatabase(c *Configs) *Database {
	return &Database{
		cfg:       c,
		fileNames: make(map[string]string),
	}
}

func (d *Database) saveToFile(base string) error {
	var err error

	if base == "relay" {
		err = d.cfg.SaveToFile(&d.relays, d.fileNames[base])
	}

	return err
}

func (d *Database) SetFileName(db string, fileName string) {
	d.fileNames[db] = fileName
}

func (d *Database) SetDBType(typ int) {
	d.dbType = typ
}

func (d *Database) LoadValues() {
	if d.dbType == DbTextType {
		// Load Relays
		d.cfg.LoadFromFile(&d.relays, d.fileNames["relay"])
	}
}

func (d *Database) Relay(name string) bool {
	for _, relay := range d.relays.Relays {
		if relay.Name == name {
			return relay.Status
		}
	}
	return false
}

func (d *Database) UpdateRelay(name string, status bool) error {
	var found bool

	for i, relay := range d.relays.Relays {
		if relay.Name == name {
			d.relays.Relays[i].Status = status
			found = true
			break
		}
	}

	if !found {
		d.relays.Relays = append(d.relays.Relays, SingleRelayDb{Name: name, Status: status})
	}

	if d.dbType == DbTextType {
		return d.saveToFile("relay")
	}

	return nil
}
