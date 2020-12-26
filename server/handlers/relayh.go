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

type RelayHandler struct {
	storage *core.Storage
	aut     *auth.Authorization
}

func NewRelayHandler(s *core.Storage, a *auth.Authorization) *RelayHandler {
	return &RelayHandler{
		storage: s,
		aut:     a,
	}
}

func (r *RelayHandler) Switch(ctx *fasthttp.RequestCtx) {
	// Find device in storage
	var device = r.storage.Device(ctx.UserValue("id").(string))
	if device == nil {
		r.Response(ctx, "switch", false, "Relay not found", nil)
		return
	}

	// Check user rights
	var _, write = r.aut.Validation(ctx.UserValue("user").(string), device.Name())
	if !write {
		r.Response(ctx, "switch", false, "Authorization failed", nil)
		return
	}

	// Process operation
	var relay = device.(*base.Relay)
	relay.Switch()

	// Send response
	r.Response(ctx, "switch", true, "", relay)
}

func (r *RelayHandler) SetStatus(ctx *fasthttp.RequestCtx) {
	// Find device in storage
	var device = r.storage.Device(ctx.UserValue("id").(string))
	if device == nil {
		r.Response(ctx, "set", false, "Relay not found", nil)
		return
	}

	// Check user rights
	var _, write = r.aut.Validation(ctx.UserValue("user").(string), device.Name())
	if !write {
		r.Response(ctx, "set", false, "Authorization failed", nil)
		return
	}

	// Process operation
	var relay = device.(*base.Relay)
	var status, err = strconv.ParseBool(ctx.UserValue("status").(string))
	if err != nil {
		r.Response(ctx, "set", false, "Fail to convert status", relay)
		return
	}
	relay.SetStatus(status)

	// Send response
	r.Response(ctx, "set", true, "", relay)
}

func (r *RelayHandler) Status(ctx *fasthttp.RequestCtx) {
	// Find device in storage
	var device = r.storage.Device(ctx.UserValue("id").(string))
	if device == nil {
		r.Response(ctx, "status", true, "Relay not found", nil)
		return
	}

	// Check user rights
	var read, _ = r.aut.Validation(ctx.UserValue("user").(string), device.Name())
	if !read {
		r.Response(ctx, "status", false, "Authorization failed", nil)
		return
	}

	// Process operation
	var relay = device.(*base.Relay)

	// Send response
	r.Response(ctx, "status", true, "", relay)
}

func (r *RelayHandler) Devices(ctx *fasthttp.RequestCtx) {
	// Find all devices
	var relays []*base.Relay
	for _, device := range r.storage.DevicesByType("relay") {
		var read, _ = r.aut.Validation(ctx.UserValue("user").(string), device.Name())
		if read {
			relays = append(relays, device.(*base.Relay))
		}
	}

	// Send response
	r.ResponseDevices(ctx, "devices", true, "", relays)
}

func (r *RelayHandler) Update(ctx *fasthttp.RequestCtx) {
	// Find device in storage
	var device = r.storage.Device(ctx.UserValue("id").(string))
	if device == nil {
		r.Response(ctx, "update", false, "Relay not found", nil)
		return
	}

	// Check user rights
	var _, write = r.aut.Validation(ctx.UserValue("user").(string), device.Name())
	if !write {
		r.Response(ctx, "switch", false, "Authorization failed", nil)
		return
	}

	// Process operation
	var relay = device.(*base.Relay)
	var state, err = strconv.ParseBool(ctx.UserValue("state").(string))
	if err != nil {
		r.Response(ctx, "update", false, "Fail to convert state", relay)
		return
	}
	relay.Update(state)

	// Send response
	r.Response(ctx, "update", true, "", relay)
}

func (r *RelayHandler) Response(ctx *fasthttp.RequestCtx, oper string, result bool, err string, relay *base.Relay) {
	ctx.Response.Header.SetContentType("application/json")

	var resp = api.RelayResponse{
		Operation: oper,
		Result:    result,
		Error:     err,
	}

	if relay != nil {
		resp.Status = relay.Status()
		resp.State = relay.State()
	}

	var bytes, _ = json.Marshal(resp)
	ctx.Write(bytes)
}

func (r *RelayHandler) ResponseDevices(ctx *fasthttp.RequestCtx, oper string, result bool, err string, relays []*base.Relay) {
	ctx.Response.Header.SetContentType("application/json")

	var relResp = api.RelayDevResponse{
		Operation: oper,
		Result:    result,
		Error:     err,
	}

	for _, relay := range relays {
		relResp.Relays = append(relResp.Relays, api.RelaySingleDevResponse{
			Name:        relay.Name(),
			Description: relay.Description(),
			Online:      relay.Online(),
			Status:      relay.Status(),
			State:       relay.State(),
		})
	}

	var bytes, _ = json.Marshal(relResp)

	ctx.Write(bytes)
}
