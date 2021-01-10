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
	"github.com/futcity/controller/server/api"
	"github.com/futcity/controller/utils"
	"github.com/valyala/fasthttp"
)

type ProfileHandler struct {
	aut     *auth.Authorization
	storage *core.Storage
	cfg     *utils.Configs
}

func NewProfileHandler(s *core.Storage, a *auth.Authorization, cfg *utils.Configs) *ProfileHandler {
	return &ProfileHandler{
		aut:     a,
		storage: s,
		cfg:     cfg,
	}
}

func (d *ProfileHandler) RemoveProfile(ctx *fasthttp.RequestCtx) {
	// Check user rights
	var admin = d.aut.IsAdmin(ctx.UserValue("user").(string))
	if !admin {
		d.response(ctx, "profile", false, "Authorization failed")
		return
	}

	// Delete profile
	var err = d.aut.DeleteProfile(ctx.UserValue("name").(string))
	if err != nil {
		d.response(ctx, "profile", false, err.Error())
		return
	}

	// Send response
	d.response(ctx, "profile", true, "")
}

func (d *ProfileHandler) AddProfile(ctx *fasthttp.RequestCtx) {
	// Check user rights
	var admin = d.aut.IsAdmin(ctx.UserValue("user").(string))
	if !admin {
		d.response(ctx, "profile", false, "Authorization failed")
		return
	}

	// Add profile
	var adm, _ = strconv.ParseBool(ctx.UserValue("admin").(string))
	d.aut.AddProfile(auth.NewProfile(ctx.UserValue("name").(string), ctx.UserValue("key").(string), adm))

	// Save to file
	d.saveProfiles(ctx)

	// Send response
	d.response(ctx, "profile", true, "")
}

func (d *ProfileHandler) AddProfileDevice(ctx *fasthttp.RequestCtx) {
	// Check user rights
	var admin = d.aut.IsAdmin(ctx.UserValue("user").(string))
	if !admin {
		d.response(ctx, "profile", false, "Authorization failed")
		return
	}

	// Add profile device
	var profile = d.aut.Profile(ctx.UserValue("name").(string))
	if profile == nil {
		d.response(ctx, "profile", false, "Profile not found")
		return
	}

	var read, _ = strconv.ParseBool(ctx.UserValue("read").(string))
	var write, _ = strconv.ParseBool(ctx.UserValue("write").(string))
	profile.AddDevice(auth.NewProfileDevice(ctx.UserValue("device").(string), read, write))

	// Save to file
	d.saveProfiles(ctx)

	// Send response
	d.response(ctx, "profile", true, "")
}

func (d *ProfileHandler) AddProfileGroup(ctx *fasthttp.RequestCtx) {
	// Check user rights
	var admin = d.aut.IsAdmin(ctx.UserValue("user").(string))
	if !admin {
		d.response(ctx, "profile", false, "Authorization failed")
		return
	}

	// Add profile device
	var profile = d.aut.Profile(ctx.UserValue("name").(string))
	if profile == nil {
		d.response(ctx, "profile", false, "Profile not found")
		return
	}

	var grp, _ = url.QueryUnescape(ctx.UserValue("group").(string))
	profile.AddGroup(grp)

	// Save to file
	d.saveProfiles(ctx)

	// Send response
	d.response(ctx, "profile", true, "")
}

func (d *ProfileHandler) RemoveProfileGroup(ctx *fasthttp.RequestCtx) {
	// Check user rights
	var admin = d.aut.IsAdmin(ctx.UserValue("user").(string))
	if !admin {
		d.response(ctx, "profile", false, "Authorization failed")
		return
	}

	// Delete profile group
	var profile = d.aut.Profile(ctx.UserValue("name").(string))
	if profile == nil {
		d.response(ctx, "profile", false, "Profile not found")
		return
	}

	var grp, _ = url.QueryUnescape(ctx.UserValue("group").(string))
	profile.RemoveGroup(grp)

	// Save to file
	d.saveProfiles(ctx)

	// Send response
	d.response(ctx, "profile", true, "")
}

func (d *ProfileHandler) RemoveProfileDevice(ctx *fasthttp.RequestCtx) {
	// Check user rights
	var admin = d.aut.IsAdmin(ctx.UserValue("user").(string))
	if !admin {
		d.response(ctx, "profile", false, "Authorization failed")
		return
	}

	// Delete profile group
	var profile = d.aut.Profile(ctx.UserValue("name").(string))
	if profile == nil {
		d.response(ctx, "profile", false, "Profile not found")
		return
	}

	profile.RemoveDevice(ctx.UserValue("device").(string))

	// Save to file
	d.saveProfiles(ctx)

	// Send response
	d.response(ctx, "profile", true, "")
}

func (d *ProfileHandler) ProfileList(ctx *fasthttp.RequestCtx) {
	// Check user rights
	var admin = d.aut.IsAdmin(ctx.UserValue("user").(string))
	if !admin {
		d.response(ctx, "profile", false, "Authorization failed")
		return
	}

	// Send response
	d.responseList(ctx, "profile", true, "")
}

func (d *ProfileHandler) saveProfiles(ctx *fasthttp.RequestCtx) {
	// Save new profile list
	var profiles configs.ProfCfg
	for _, profile := range d.aut.Profiles() {
		var p = configs.ProfileCfg{
			Name:  profile.Name(),
			Key:   profile.APIKey(),
			Admin: profile.Admin(),
		}

		p.Groups = make([]string, len(*profile.Groups()))
		copy(p.Groups, *profile.Groups())

		for _, dev := range profile.Devices() {
			p.Devices = append(p.Devices, configs.ProfDevCfg{
				Name:  dev.Name(),
				Read:  dev.Read(),
				Write: dev.Write(),
			})
		}
		profiles.Profiles = append(profiles.Profiles, p)
	}
	var err = d.cfg.SaveToFile(&profiles, "./profiles.conf")
	if err != nil {
		d.response(ctx, "profile", false, err.Error())
		return
	}
}

func (d *ProfileHandler) response(ctx *fasthttp.RequestCtx, oper string, result bool, err string) {
	ctx.Response.Header.SetContentType("application/json")

	var devResp = api.ProfileListResponse{
		Operation: oper,
		Result:    result,
		Error:     err,
	}

	var bytes, _ = json.Marshal(devResp)

	ctx.Write(bytes)
}

func (d *ProfileHandler) responseList(ctx *fasthttp.RequestCtx, oper string, result bool, err string) {
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

	var bytes, _ = json.Marshal(devResp)

	ctx.Write(bytes)
}
