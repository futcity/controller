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
	"os"

	"github.com/futcity/controller/auth"
	"github.com/futcity/controller/configs"
	"github.com/futcity/controller/core"
	"github.com/futcity/controller/db"
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
	db      *db.Database
}

// NewApp Make new struct
func NewApp(s *core.Storage, srv *server.WebServer, a *auth.Authorization, l *utils.Log,
	c *utils.Configs, d *db.Database) *App {
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
	var ac configs.AppCfg

	if len(os.Args) > 1 {
		var err = a.cfg.LoadFromFile(&ac, os.Args[1])
		if err != nil {
			a.log.Error("APP", "Fail to load configs", err.Error())
		}
	} else {
		a.log.Error("APP", "Fail to load configs", "Error args count. Please add configs file path")
		return
	}
	a.log.Info("APP", "Configs was loaded")

	//
	// Load database
	//
	a.db.SetDBType(ac.Db.Type)
	if ac.Db.Type == "text" {
		for _, file := range ac.Db.Files {
			a.db.AddFilename(file.Name, file.Path)
		}
	}

	var err = a.db.LoadDeviceBase()
	if err != nil {
		a.log.Error("APP", "Fail to load device database", err.Error())
		return
	}
	err = a.db.LoadProfileBase()
	if err != nil {
		a.log.Error("APP", "Fail to load profile database", err.Error())
		return
	}
	err = a.db.LoadRelayBase()
	if err != nil {
		a.log.Error("APP", "Fail to load relay database", err.Error())
		return
	}

	//
	// Starting server
	//
	a.log.Info("APP", "Starting server...")
	err = a.server.Start(ac.Server.IP, ac.Server.Port)
	if err != nil {
		a.log.Error("APP", "Fail to start web server", err.Error())
	}
}
