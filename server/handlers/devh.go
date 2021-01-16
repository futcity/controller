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
	"net/url"
	"strconv"

	"github.com/futcity/controller/auth"
	"github.com/futcity/controller/core"
	"github.com/futcity/controller/core/devices"
	"github.com/futcity/controller/db"
	"github.com/futcity/controller/server/api"
	"github.com/futcity/controller/utils"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

type DeviceHandler struct {
	aut     *auth.Authorization
	storage *core.Storage
	db      *db.Database
	log     *utils.Log
}

func NewDeviceHandler(s *core.Storage, a *auth.Authorization, db *db.Database,
	l *utils.Log) *DeviceHandler {
	return &DeviceHandler{
		aut:     a,
		storage: s,
		db:      db,
		log:     l,
	}
}

func (d *DeviceHandler) DeviceByDescription(ctx *fasthttp.RequestCtx) {
	// Find device in storage
	var desc, _ = url.QueryUnescape(ctx.UserValue("desc").(string))
	var device, err = d.storage.DeviceByDescription(desc)
	if err != nil {
		d.response(ctx, "Get device by desc", false, err.Error(), "")
		return
	}

	// Check user rights
	var _, write = d.aut.Validation(ctx.UserValue("user").(string), device.Name())
	if !write {
		d.response(ctx, "Get device by desc", false, "Authorization failed", "")
		return
	}

	// Send response
	d.response(ctx, "Get device by desc", true, "", device.Name())
}

func (d *DeviceHandler) RemoveDevice(ctx *fasthttp.RequestCtx) {
	// Find device in storage
	var id, _ = strconv.Atoi(ctx.UserValue("id").(string))

	var device, err = d.storage.DeviceByID(id)
	if err != nil {
		d.response(ctx, "Remove device", false, err.Error(), "")
		return
	}

	// Check user rights
	var _, write = d.aut.Validation(ctx.UserValue("user").(string), device.Name())
	if !write {
		d.response(ctx, "Remove device", false, "Authorization failed", "")
		return
	}

	// Delete device
	err = d.storage.RemoveByID(id)
	if err != nil {
		d.response(ctx, "Remove device", false, err.Error(), "")
		return
	}

	// Send response
	d.response(ctx, "Remove device", true, "", device.Name())
}

func (d *DeviceHandler) AddDevice(ctx *fasthttp.RequestCtx) {
	// Check user rights
	var admin = d.aut.IsAdmin(ctx.UserValue("user").(string))
	if !admin {
		d.response(ctx, "Add device", false, "Authorization failed", "")
		return
	}

	// Add device to storage
	var desc, _ = url.QueryUnescape(ctx.UserValue("desc").(string))
	d.storage.AddDevice(ctx.UserValue("name").(string), desc, ctx.UserValue("type").(string))

	// Save new devices list
	var err = d.db.SaveDeviceBase()
	if err != nil {
		d.response(ctx, "Save device", false, err.Error(), ctx.UserValue("name").(string))
		return
	}

	// Send response
	d.response(ctx, "Add device", true, "", ctx.UserValue("name").(string))
}

func (d *DeviceHandler) DeviceList(ctx *fasthttp.RequestCtx) {
	var devices []devices.IDevice

	// Check user rights and add device to list
	for _, device := range d.storage.Devices() {
		var _, write = d.aut.Validation(ctx.UserValue("user").(string), device.Name())
		if write {
			devices = append(devices, device)
		}
	}

	// Send response
	d.responseList(ctx, "Devices list", true, "", devices)
}

func (d *DeviceHandler) response(ctx *fasthttp.RequestCtx, oper string, result bool, err string, name string) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	ctx.Response.Header.SetContentType("application/json")

	var devResp = api.DeviceResponse{
		Operation: oper,
		Result:    result,
		Error:     err,
		Name:      name,
	}

	if result {
		d.log.Info("GROUPH", oper)
	} else {
		d.log.Error("GROUPH", oper, err)
	}

	var bytes, _ = json.Marshal(devResp)

	ctx.Write(bytes)
}

func (d *DeviceHandler) responseList(ctx *fasthttp.RequestCtx, oper string, result bool, err string, devices []devices.IDevice) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	ctx.Response.Header.SetContentType("application/json")

	var devResp = api.DeviceListResponse{
		Operation: oper,
		Result:    result,
		Error:     err,
	}

	for _, device := range devices {
		devResp.Devices = append(devResp.Devices, api.DeviceSingleResponse{
			ID:          device.ID(),
			Name:        device.Name(),
			Description: device.Description(),
			Type:        device.Type(),
			Online:      device.Online(),
		})
	}

	if result {
		d.log.Info("GROUPH", oper)
	} else {
		d.log.Error("GROUPH", oper, err)
	}

	var bytes, _ = json.Marshal(devResp)

	ctx.Write(bytes)
}
