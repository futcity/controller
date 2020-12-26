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
	"fmt"

	"github.com/futcity/controller/auth"
	"github.com/futcity/controller/core"
	"github.com/futcity/controller/server"
	"github.com/futcity/controller/server/handlers"
	"github.com/futcity/controller/utils"
	"go.uber.org/dig"
)

func main() {
	container := dig.New()

	container.Provide(utils.NewLog)
	container.Provide(utils.NewConfigs)
	container.Provide(utils.NewDatabase)

	container.Provide(auth.NewAuthorization)
	container.Provide(core.NewStorage)

	container.Provide(handlers.NewGroupHandler)
	container.Provide(handlers.NewRelayHandler)
	container.Provide(server.NewWebServer)

	container.Provide(NewApp)

	err := container.Invoke(func(app *App) {
		app.Start()
	})

	if err != nil {
		fmt.Printf("Fatal error: %s\n", err.Error())
	}
}
