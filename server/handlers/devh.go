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
	"net/url"
	"strconv"

	"github.com/futcity/controller/auth"
	"github.com/futcity/controller/configs"
	"github.com/futcity/controller/core"
	"github.com/futcity/controller/core/devices"
	"github.com/futcity/controller/server/api"
	"github.com/futcity/controller/utils"
	"github.com/valyala/fasthttp"
)

type DeviceHandler struct {
	aut     *auth.Authorization
	storage *core.Storage
	cfg     *utils.Configs
}

func NewDeviceHandler(s *core.Storage, a *auth.Authorization, cfg *utils.Configs) *DeviceHandler {
	return &DeviceHandler{
		aut:     a,
		storage: s,
		cfg:     cfg,
	}
}

func (d *DeviceHandler) DeviceByDescription(ctx *fasthttp.RequestCtx) {
	// Find device in storage
	var desc, _ = url.QueryUnescape(ctx.UserValue("desc").(string))
	var device, err = d.storage.DeviceByDescription(desc)
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

func (d *DeviceHandler) RemoveDevice(ctx *fasthttp.RequestCtx) {
	// Find device in storage
	var id, _ = strconv.Atoi(ctx.UserValue("id").(string))

	var device, err = d.storage.DeviceByID(id)
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

	// Delete device
	err = d.storage.RemoveByID(id)
	if err != nil {
		d.Response(ctx, "device", false, err.Error(), "")
		return
	}

	// Send response
	d.Response(ctx, "device", true, "", device.Name())
}

func (d *DeviceHandler) AddDevice(ctx *fasthttp.RequestCtx) {
	// Check user rights
	var admin = d.aut.IsAdmin(ctx.UserValue("user").(string))
	if !admin {
		d.Response(ctx, "device", false, "Authorization failed", "")
		return
	}

	// Add device to storage
	var desc, _ = url.QueryUnescape(ctx.UserValue("desc").(string))
	d.storage.AddDevice(ctx.UserValue("name").(string), desc, ctx.UserValue("type").(string))

	// Save new devices list
	var devices configs.DevCfg
	for _, device := range d.storage.Devices() {
		devices.Devices = append(devices.Devices, configs.DeviceCfg{
			Name:        device.Name(),
			Description: device.Description(),
			Type:        device.Type(),
		})
	}
	var err = d.cfg.SaveToFile(&devices, "./devices.conf")
	if err != nil {
		d.Response(ctx, "device", true, "", ctx.UserValue("name").(string))
		return
	}

	// Send response
	d.Response(ctx, "device", true, "", ctx.UserValue("name").(string))
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
	d.ResponseList(ctx, "device", true, "", devices)
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

func (d *DeviceHandler) ResponseList(ctx *fasthttp.RequestCtx, oper string, result bool, err string, devices []devices.IDevice) {
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

	var bytes, _ = json.Marshal(devResp)

	ctx.Write(bytes)
}
