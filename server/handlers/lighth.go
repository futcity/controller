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
	"strconv"

	"github.com/futcity/controller/auth"
	"github.com/futcity/controller/core"
	"github.com/futcity/controller/core/devices/base"
	"github.com/futcity/controller/server/api"
	"github.com/valyala/fasthttp"
)

type LightHandler struct {
	storage *core.Storage
	aut     *auth.Authorization
}

func NewLightHandler(s *core.Storage, a *auth.Authorization) *LightHandler {
	return &LightHandler{
		storage: s,
		aut:     a,
	}
}

func (l *LightHandler) Switch(ctx *fasthttp.RequestCtx) {
	// Find device in storage
	var device = l.storage.Device(ctx.UserValue("id").(string))
	if device == nil {
		l.Response(ctx, "switch", false, "Light switch not found", nil)
		return
	}

	// Check user rights
	var _, write = l.aut.Validation(ctx.UserValue("user").(string), device.Name())
	if !write {
		l.Response(ctx, "switch", false, "Authorization failed", nil)
		return
	}

	// Process operation
	var light = device.(*base.Light)
	light.Switch()

	// Send response
	l.Response(ctx, "switch", true, "", light)
}

func (l *LightHandler) SetStatus(ctx *fasthttp.RequestCtx) {
	// Find device in storage
	var device = l.storage.Device(ctx.UserValue("id").(string))
	if device == nil {
		l.Response(ctx, "set", false, "Light switch not found", nil)
		return
	}

	// Check user rights
	var _, write = l.aut.Validation(ctx.UserValue("user").(string), device.Name())
	if !write {
		l.Response(ctx, "set", false, "Authorization failed", nil)
		return
	}

	// Process operation
	var light = device.(*base.Light)
	var status, err = strconv.ParseBool(ctx.UserValue("status").(string))
	if err != nil {
		l.Response(ctx, "set", false, "Fail to convert status", light)
		return
	}
	light.SetStatus(status)

	// Send response
	l.Response(ctx, "set", true, "", light)
}

func (l *LightHandler) Status(ctx *fasthttp.RequestCtx) {
	// Find device in storage
	var device = l.storage.Device(ctx.UserValue("id").(string))
	if device == nil {
		l.Response(ctx, "status", true, "Light switch not found", nil)
		return
	}

	// Check user rights
	var read, _ = l.aut.Validation(ctx.UserValue("user").(string), device.Name())
	if !read {
		l.Response(ctx, "status", false, "Authorization failed", nil)
		return
	}

	// Process operation
	var light = device.(*base.Light)

	// Send response
	l.Response(ctx, "status", true, "", light)
}

func (l *LightHandler) Devices(ctx *fasthttp.RequestCtx) {
	// Find all devices
	var relays []*base.Light
	for _, device := range l.storage.DevicesByType("light") {
		var read, _ = l.aut.Validation(ctx.UserValue("user").(string), device.Name())
		if read {
			relays = append(relays, device.(*base.Light))
		}
	}

	// Send response
	l.ResponseDevices(ctx, "devices", true, "", relays)
}

func (l *LightHandler) Update(ctx *fasthttp.RequestCtx) {
	// Find device in storage
	var device = l.storage.Device(ctx.UserValue("id").(string))
	if device == nil {
		l.Response(ctx, "update", false, "Light switch not found", nil)
		return
	}

	// Check user rights
	var _, write = l.aut.Validation(ctx.UserValue("user").(string), device.Name())
	if !write {
		l.Response(ctx, "switch", false, "Authorization failed", nil)
		return
	}

	// Process operation
	var light = device.(*base.Light)
	var state, err = strconv.ParseBool(ctx.UserValue("state").(string))
	if err != nil {
		l.Response(ctx, "update", false, "Fail to convert state", light)
		return
	}
	light.Update(state)

	// Send response
	l.Response(ctx, "update", true, "", light)
}

func (l *LightHandler) Response(ctx *fasthttp.RequestCtx, oper string, result bool, err string, light *base.Light) {
	ctx.Response.Header.SetContentType("application/json")

	var bytes, _ = json.Marshal(api.RelayResponse{
		Operation: oper,
		Result:    result,
		Error:     err,
		Status:    light.Status(),
		State:     light.State(),
	})

	ctx.Write(bytes)
}

func (l *LightHandler) ResponseDevices(ctx *fasthttp.RequestCtx, oper string, result bool, err string, lights []*base.Light) {
	ctx.Response.Header.SetContentType("application/json")

	var lghtResp = api.LightDevResponse{
		Operation: oper,
		Result:    result,
		Error:     err,
	}

	for _, light := range lights {
		lghtResp.Lights = append(lghtResp.Lights, api.LightSingleDevResponse{
			Name:        light.Name(),
			Description: light.Description(),
			Online:      light.Online(),
			Status:      light.Status(),
			State:       light.State(),
		})
	}

	var bytes, _ = json.Marshal(lghtResp)

	ctx.Write(bytes)
}
