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
	"github.com/futcity/controller/db"
	"github.com/futcity/controller/server/api"
	"github.com/futcity/controller/utils"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

type ProfileHandler struct {
	aut     *auth.Authorization
	storage *core.Storage
	db      *db.Database
	log     *utils.Log
}

func NewProfileHandler(s *core.Storage, a *auth.Authorization, db *db.Database,
	l *utils.Log) *ProfileHandler {
	return &ProfileHandler{
		aut:     a,
		storage: s,
		db:      db,
		log:     l,
	}
}

func (d *ProfileHandler) RemoveProfile(ctx *fasthttp.RequestCtx) {
	// Check user rights
	var admin = d.aut.IsAdmin(ctx.UserValue("user").(string))
	if !admin {
		d.response(ctx, "Remove profile", false, "Authorization failed")
		return
	}

	// Delete profile
	var err = d.aut.DeleteProfile(ctx.UserValue("name").(string))
	if err != nil {
		d.response(ctx, "Remove profile", false, err.Error())
		return
	}

	// Save new profile list
	err = d.db.SaveProfileBase()
	if err != nil {
		d.response(ctx, "Save profile", false, err.Error())
		return
	}

	// Send response
	d.response(ctx, "Remove profile", true, "")
}

func (d *ProfileHandler) AddProfile(ctx *fasthttp.RequestCtx) {
	// Check user rights
	var admin = d.aut.IsAdmin(ctx.UserValue("user").(string))
	if !admin {
		d.response(ctx, "Add profile", false, "Authorization failed")
		return
	}

	// Add profile
	var adm, _ = strconv.ParseBool(ctx.UserValue("admin").(string))
	d.aut.AddProfile(auth.NewProfile(ctx.UserValue("name").(string), ctx.UserValue("key").(string), adm))

	// Save new profile list
	var err = d.db.SaveProfileBase()
	if err != nil {
		d.response(ctx, "Add profile", false, err.Error())
		return
	}

	// Send response
	d.response(ctx, "Add profile", true, "")
}

func (d *ProfileHandler) AddProfileDevice(ctx *fasthttp.RequestCtx) {
	// Check user rights
	var admin = d.aut.IsAdmin(ctx.UserValue("user").(string))
	if !admin {
		d.response(ctx, "Add profile device", false, "Authorization failed")
		return
	}

	// Add profile device
	var profile = d.aut.Profile(ctx.UserValue("name").(string))
	if profile == nil {
		d.response(ctx, "Add profile device", false, "Profile not found")
		return
	}

	var read, _ = strconv.ParseBool(ctx.UserValue("read").(string))
	var write, _ = strconv.ParseBool(ctx.UserValue("write").(string))
	profile.AddDevice(auth.NewProfileDevice(ctx.UserValue("device").(string), read, write))

	// Save new profile list
	var err = d.db.SaveProfileBase()
	if err != nil {
		d.response(ctx, "Add profile device", false, err.Error())
		return
	}

	// Send response
	d.response(ctx, "Add profile device", true, "")
}

func (d *ProfileHandler) AddProfileGroup(ctx *fasthttp.RequestCtx) {
	// Check user rights
	var admin = d.aut.IsAdmin(ctx.UserValue("user").(string))
	if !admin {
		d.response(ctx, "Add profile group", false, "Authorization failed")
		return
	}

	// Add profile device
	var profile = d.aut.Profile(ctx.UserValue("name").(string))
	if profile == nil {
		d.response(ctx, "Add profile group", false, "Profile not found")
		return
	}

	var grp, _ = url.QueryUnescape(ctx.UserValue("group").(string))
	profile.AddGroup(grp)

	// Save new profile list
	var err = d.db.SaveProfileBase()
	if err != nil {
		d.response(ctx, "Add profile group", false, err.Error())
		return
	}

	// Send response
	d.response(ctx, "Add profile group", true, "")
}

func (d *ProfileHandler) RemoveProfileGroup(ctx *fasthttp.RequestCtx) {
	// Check user rights
	var admin = d.aut.IsAdmin(ctx.UserValue("user").(string))
	if !admin {
		d.response(ctx, "Remove profile group", false, "Authorization failed")
		return
	}

	// Delete profile group
	var profile = d.aut.Profile(ctx.UserValue("name").(string))
	if profile == nil {
		d.response(ctx, "Remove profile group", false, "Profile not found")
		return
	}

	var grp, _ = url.QueryUnescape(ctx.UserValue("group").(string))
	profile.RemoveGroup(grp)

	// Save new profile list
	var err = d.db.SaveProfileBase()
	if err != nil {
		d.response(ctx, "Remove profile group", false, err.Error())
		return
	}

	// Send response
	d.response(ctx, "Remove profile group", true, "")
}

func (d *ProfileHandler) RemoveProfileDevice(ctx *fasthttp.RequestCtx) {
	// Check user rights
	var admin = d.aut.IsAdmin(ctx.UserValue("user").(string))
	if !admin {
		d.response(ctx, "Remove profile device", false, "Authorization failed")
		return
	}

	// Delete profile group
	var profile = d.aut.Profile(ctx.UserValue("name").(string))
	if profile == nil {
		d.response(ctx, "Remove profile device", false, "Profile not found")
		return
	}

	profile.RemoveDevice(ctx.UserValue("device").(string))

	// Save new profile list
	var err = d.db.SaveProfileBase()
	if err != nil {
		d.response(ctx, "Remove profile device", false, err.Error())
		return
	}

	// Send response
	d.response(ctx, "Remove profile device", true, "")
}

func (d *ProfileHandler) ProfileList(ctx *fasthttp.RequestCtx) {
	// Check user rights
	var admin = d.aut.IsAdmin(ctx.UserValue("user").(string))
	if !admin {
		d.response(ctx, "Profile list", false, "Authorization failed")
		return
	}

	// Send response
	d.responseList(ctx, "Profile list", true, "")
}

func (d *ProfileHandler) response(ctx *fasthttp.RequestCtx, oper string, result bool, err string) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	ctx.Response.Header.SetContentType("application/json")

	var devResp = api.ProfileListResponse{
		Operation: oper,
		Result:    result,
		Error:     err,
	}

	if result {
		d.log.Info("PROFILEH", oper)
	} else {
		d.log.Error("PROFILEH", oper, err)
	}

	var bytes, _ = json.Marshal(devResp)

	ctx.Write(bytes)
}

func (d *ProfileHandler) responseList(ctx *fasthttp.RequestCtx, oper string, result bool, err string) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	ctx.Response.Header.SetContentType("application/json")

	var devResp = api.ProfileListResponse{
		Operation: oper,
		Result:    result,
		Error:     err,
	}

	if result {
		for _, prof := range d.aut.Profiles() {
			var p = api.ProfileSingleResponse{
				Name:  prof.Name(),
				Key:   prof.APIKey(),
				Admin: prof.Admin(),
			}
			p.Groups = make([]string, len(*prof.Groups()))
			copy(p.Groups, *prof.Groups())
			for _, dev := range prof.Devices() {
				p.Devices = append(p.Devices, api.ProfileDeviceResponse{
					Name:  dev.Name(),
					Read:  dev.Read(),
					Write: dev.Write(),
				})
			}
			devResp.Profiles = append(devResp.Profiles, p)
		}
	}

	if result {
		d.log.Info("PROFILEH", oper)
	} else {
		d.log.Error("PROFILEH", oper, err)
	}

	var bytes, _ = json.Marshal(devResp)

	ctx.Write(bytes)
}
