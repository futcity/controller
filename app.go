///////////////////////////////////////////////////////////////////
//
// Future City Project
//
// Copyright (C) 2020-2021 Sergey Denisov. GPLv3
//
// Written by Sergey Denisov aka LittleBuster (DenisovS21@gmail.com)
//
///////////////////////////////////////////////////////////////////

package main

import (
	"github.com/futcity/controller/auth"
	"github.com/futcity/controller/configs"
	"github.com/futcity/controller/core"
	"github.com/futcity/controller/server"
	"github.com/futcity/controller/utils"
)

// App Application main module
type App struct {
	storage *core.Storage
	server  *server.WebServer
	aut     *auth.Authorization
	log     *utils.Log
	cfg     *utils.Configs
	db      *utils.Database
}

// NewApp Make new struct
func NewApp(s *core.Storage, srv *server.WebServer, a *auth.Authorization, l *utils.Log,
	c *utils.Configs, d *utils.Database) *App {
	return &App{
		storage: s,
		server:  srv,
		aut:     a,
		log:     l,
		cfg:     c,
		db:      d,
	}
}

// Start applications
func (a *App) Start() {
	//
	// Init logger
	//
	a.log.SetPath("./")

	//
	// Loading configs
	//
	var dc = &configs.DevCfg{}
	var pc = &configs.ProfCfg{}
	var ac = &configs.AppCfg{}

	var err = a.cfg.LoadFromFile(dc, "devices.conf")
	if err != nil {
		a.log.Error("APP", "Fail to load device configs", err.Error())
		return
	}
	err = a.cfg.LoadFromFile(pc, "profiles.conf")
	if err != nil {
		a.log.Error("APP", "Fail to load profiles configs", err.Error())
		return
	}
	err = a.cfg.LoadFromFile(ac, "fc.conf")
	if err != nil {
		a.log.Error("APP", "Fail to load app configs", err.Error())
		return
	}

	//
	// Load database
	//
	if ac.Db.Type == "text" {
		for _, dbFile := range ac.Db.Files {
			a.db.SetFileName(dbFile.Name, dbFile.Path)
			a.log.Info("APP", "Add new db filename \""+dbFile.Name+"\"")
		}
	}

	a.db.LoadValues()

	//
	// Applying configs
	//

	// Create devices
	for _, dev := range dc.Devices {
		a.storage.AddDevice(dev.Name, dev.Description, dev.Type)
		a.log.Info("APP", "Add new device \""+dev.Name+"\" desc \""+dev.Description+"\" type \""+dev.Type+"\"")
	}

	// Create user profiles
	for _, prof := range pc.Profiles {
		a.log.Info("APP", "Add new profile name \""+prof.Name+"\"")

		var profile = auth.NewProfile(prof.Name, prof.Key, prof.Admin)
		for _, pdev := range prof.Devices {
			profile.AddDevice(auth.NewProfileDevice(pdev.Name, pdev.Read, pdev.Write))
			a.log.Info("APP", "Add new profile \""+prof.Name+"\" device \""+pdev.Name+"\"")
		}
		for _, grp := range prof.Groups {
			profile.AddGroup(grp)
			a.log.Info("APP", "Add new profile \""+prof.Name+"\" group \""+grp+"\"")
		}

		a.aut.AddProfile(profile)
	}

	//
	// Free configs
	//
	dc = nil
	pc = nil

	//
	// Starting server
	//
	a.log.Info("APP", "Starting server...")
	err = a.server.Start(ac.Server.IP, ac.Server.Port)
	if err != nil {
		a.log.Error("APP", "Fail to start web server", err.Error())
	}
}
