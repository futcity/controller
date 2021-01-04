///////////////////////////////////////////////////////////////////
//
// Future City Project
//
// Copyright (C) 2020-2021 Sergey Denisov. GPLv3
//
// Written by Sergey Denisov aka LittleBuster (DenisovS21@gmail.com)
//
///////////////////////////////////////////////////////////////////

package handlers

import (
	"encoding/json"

	"github.com/futcity/controller/auth"
	"github.com/futcity/controller/core"
	"github.com/futcity/controller/server/api"
	"github.com/valyala/fasthttp"
)

type DeviceHandler struct {
	aut     *auth.Authorization
	storage *core.Storage
}

func NewDeviceHandler(s *core.Storage, a *auth.Authorization) *DeviceHandler {
	return &DeviceHandler{
		aut:     a,
		storage: s,
	}
}

func (d *DeviceHandler) DeviceByDescription(ctx *fasthttp.RequestCtx) {
	// Find device in storage
	var device, err = d.storage.DeviceByDescription(ctx.UserValue("desc").(string))
	if err != nil {
		d.Response(ctx, "device", false, err.Error(), "")
		return
	}

	// Check user rights
	var _, write = d.aut.Validation(ctx.UserValue("user").(string), device.Name())
	if !write {
		d.Response(ctx, "device", false, "Authorization failed", "")
		return
	}

	// Send response
	d.Response(ctx, "device", true, "", device.Name())
}

func (d *DeviceHandler) Response(ctx *fasthttp.RequestCtx, oper string, result bool, err string, name string) {
	ctx.Response.Header.SetContentType("application/json")

	var devResp = api.DeviceResponse{
		Operation: oper,
		Result:    result,
		Error:     err,
		Name:      name,
	}

	var bytes, _ = json.Marshal(devResp)

	ctx.Write(bytes)
}
