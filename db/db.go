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

import (
	"strconv"

	"github.com/futcity/controller/auth"
	"github.com/futcity/controller/core"
	"github.com/futcity/controller/core/devices/base"
	"github.com/futcity/controller/utils"
)

const (
	DbTextType = iota
)

type Database struct {
	// Dependencies
	cfg     *utils.Configs
	aut     *auth.Authorization
	storage *core.Storage
	log     *utils.Log

	// Local variables
	fileNames map[string]string
	dbType    int
}

func NewDatabase(c *utils.Configs, a *auth.Authorization, s *core.Storage,
	l *utils.Log) *Database {
	return &Database{
		cfg:       c,
		aut:       a,
		storage:   s,
		log:       l,
		fileNames: make(map[string]string),
	}
}

func (d *Database) SetDBType(typ string) {
	if typ == "text" {
		d.dbType = DbTextType
	}
}

func (d *Database) AddFilename(db string, fileName string) {
	d.fileNames[db] = fileName
}

//
// Main database managment
//

func (d *Database) LoadDeviceBase() error {
	var devices DeviceDB

	if d.dbType == DbTextType {
		var err = d.cfg.LoadFromFile(&devices, d.fileNames["device"])
		if err != nil {
			return err
		}

		for _, device := range devices.Devices {
			d.storage.AddDevice(device.Name, device.Description, device.Type)
			d.log.Info("DB", "Add new device \""+device.Name+"\" desc \""+device.Description+"\" type \""+device.Type+"\"")
		}
	}

	return nil
}

func (d *Database) SaveDeviceBase() error {
	var devices DeviceDB

	if d.dbType == DbTextType {
		for _, device := range d.storage.Devices() {
			devices.Devices = append(devices.Devices, SingleDeviceDb{
				Name:        device.Name(),
				Description: device.Description(),
				Type:        device.Type(),
			})
		}
	}

	return d.cfg.SaveToFile(&devices, d.fileNames["device"])
}

func (d *Database) LoadProfileBase() error {
	var profiles ProfileDB

	if d.dbType == DbTextType {
		var err = d.cfg.LoadFromFile(&profiles, d.fileNames["profile"])
		if err != nil {
			return err
		}

		for _, prof := range profiles.Profiles {
			d.log.Info("DB", "Add new profile name \""+prof.Name+"\"")

			var profile = auth.NewProfile(prof.Name, prof.Key, prof.Admin)
			for _, pdev := range prof.Devices {
				profile.AddDevice(auth.NewProfileDevice(pdev.Name, pdev.Read, pdev.Write))
				d.log.Info("DB", "Add new profile \""+prof.Name+"\" device \""+pdev.Name+"\"")
			}
			for _, grp := range prof.Groups {
				profile.AddGroup(grp)
				d.log.Info("DB", "Add new profile \""+prof.Name+"\" group \""+grp+"\"")
			}

			d.aut.AddProfile(profile)
		}
	}

	return nil
}

func (d *Database) SaveProfileBase() error {
	var profiles ProfileDB

	for _, profile := range d.aut.Profiles() {
		var p = SingleProfileDB{
			Name:  profile.Name(),
			Key:   profile.APIKey(),
			Admin: profile.Admin(),
		}

		p.Groups = make([]string, len(*profile.Groups()))
		copy(p.Groups, *profile.Groups())

		for _, dev := range profile.Devices() {
			p.Devices = append(p.Devices, ProfileDeivceDB{
				Name:  dev.Name(),
				Read:  dev.Read(),
				Write: dev.Write(),
			})
		}
		profiles.Profiles = append(profiles.Profiles, p)
	}

	return d.cfg.SaveToFile(&profiles, d.fileNames["profile"])
}

//
// Devices database managment
//

func (d *Database) LoadRelayBase() error {
	var relays RelaysDB

	if d.dbType == DbTextType {
		var err = d.cfg.LoadFromFile(&relays, d.fileNames["relay"])
		if err != nil {
			return err
		}

		for _, relay := range relays.Relays {
			var r = d.storage.Device(relay.Name)
			if r != nil {
				r.(*base.Relay).SetStatus(relay.Status)
				d.log.Info("DB", "Load relay status \""+relay.Name+"\" status \""+strconv.FormatBool(relay.Status)+"\"")
			}
		}
	}

	return nil
}

func (d *Database) SaveRelayBase(name string, status bool) error {
	var relays RelaysDB

	if d.dbType == DbTextType {
		for _, relay := range d.storage.DevicesByType("relay") {
			relays.Relays = append(relays.Relays, SingleRelayDB{
				Name:   relay.Name(),
				Status: relay.(*base.Relay).Status(),
			})
		}
	}

	return d.cfg.SaveToFile(&relays, d.fileNames["relay"])
}
