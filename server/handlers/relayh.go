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
	"strconv"

	"github.com/futcity/controller/auth"
	"github.com/futcity/controller/core"
	"github.com/futcity/controller/core/devices/base"
	"github.com/futcity/controller/db"
	"github.com/futcity/controller/server/api"
	"github.com/futcity/controller/utils"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

type RelayHandler struct {
	storage *core.Storage
	aut     *auth.Authorization
	db      *db.Database
	log     *utils.Log
}

func NewRelayHandler(s *core.Storage, a *auth.Authorization, l *utils.Log, db *db.Database) *RelayHandler {
	return &RelayHandler{
		storage: s,
		aut:     a,
		log:     l,
		db:      db,
	}
}

func (r *RelayHandler) Switch(ctx *fasthttp.RequestCtx) {
	// Find device in storage
	var device = r.storage.Device(ctx.UserValue("id").(string))
	if device == nil {
		r.response(ctx, "Switch relay", false, "Relay not found", nil)
		return
	}

	// Check user rights
	var _, write = r.aut.Validation(ctx.UserValue("user").(string), device.Name())
	if !write {
		r.response(ctx, "Switch relay", false, "Authorization failed", nil)
		return
	}

	// Process operation
	var relay = device.(*base.Relay)
	relay.Switch()

	// Save to database
	var err = r.db.SaveRelayBase(relay.Name(), relay.Status())
	if err != nil {
		r.response(ctx, "Save to relay database", false, err.Error(), relay)
		return
	}

	// Send response
	r.response(ctx, "Switch relay", true, "", relay)
}

func (r *RelayHandler) SetStatus(ctx *fasthttp.RequestCtx) {
	// Find device in storage
	var device = r.storage.Device(ctx.UserValue("id").(string))
	if device == nil {
		r.response(ctx, "Set relay status", false, "Relay not found", nil)
		return
	}

	// Check user rights
	var _, write = r.aut.Validation(ctx.UserValue("user").(string), device.Name())
	if !write {
		r.response(ctx, "Set relay status", false, "Authorization failed", nil)
		return
	}

	// Process operation
	var relay = device.(*base.Relay)
	var status, err = strconv.ParseBool(ctx.UserValue("status").(string))
	if err != nil {
		r.response(ctx, "Set relay status", false, "Fail to convert status", relay)
		return
	}
	relay.SetStatus(status)

	// Save to database
	err = r.db.SaveRelayBase(relay.Name(), relay.Status())
	if err != nil {
		r.response(ctx, "Save to relay database", false, err.Error(), relay)
		return
	}

	// Send response
	r.response(ctx, "Set relay status", true, "", relay)
}

func (r *RelayHandler) Status(ctx *fasthttp.RequestCtx) {
	// Find device in storage
	var device = r.storage.Device(ctx.UserValue("id").(string))
	if device == nil {
		r.response(ctx, "Get relay status", true, "Relay not found", nil)
		return
	}

	// Check user rights
	var read, _ = r.aut.Validation(ctx.UserValue("user").(string), device.Name())
	if !read {
		r.response(ctx, "Get relay status", false, "Authorization failed", nil)
		return
	}

	// Process operation
	var relay = device.(*base.Relay)

	// Send response
	r.response(ctx, "Get relay status", true, "", relay)
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
	r.responseList(ctx, "Get relays list", true, "", relays)
}

func (r *RelayHandler) Update(ctx *fasthttp.RequestCtx) {
	// Find device in storage
	var device = r.storage.Device(ctx.UserValue("id").(string))
	if device == nil {
		r.response(ctx, "Update relay", false, "Relay not found", nil)
		return
	}

	// Check user rights
	var _, write = r.aut.Validation(ctx.UserValue("user").(string), device.Name())
	if !write {
		r.response(ctx, "Update relay", false, "Authorization failed", nil)
		return
	}

	// Process operation
	var relay = device.(*base.Relay)
	var state, err = strconv.ParseBool(ctx.UserValue("state").(string))
	if err != nil {
		r.response(ctx, "Update relay", false, "Fail to convert state", relay)
		return
	}
	relay.Update(state)

	// Send response
	r.response(ctx, "Update relay", true, "", relay)
}

func (r *RelayHandler) response(ctx *fasthttp.RequestCtx, oper string, result bool, err string, relay *base.Relay) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
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

	if result && oper != "Update relay" {
		r.log.Info("RELAYH", oper)
	} else {
		r.log.Error("RELAYH", oper, err)
	}

	var bytes, _ = json.Marshal(resp)
	ctx.Write(bytes)
}

func (r *RelayHandler) responseList(ctx *fasthttp.RequestCtx, oper string, result bool, err string, relays []*base.Relay) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
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

	if result {
		r.log.Info("RELAYH", oper)
	} else {
		r.log.Error("RELAYH", oper, err)
	}

	var bytes, _ = json.Marshal(relResp)

	ctx.Write(bytes)
}
