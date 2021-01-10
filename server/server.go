///////////////////////////////////////////////////////////////////
//
// Future City Project
//
// Copyright (C) 2020-2021 Sergey Denisov. GPLv3
//
// Written by Sergey Denisov aka LittleBuster (DenisovS21@gmail.com)
//
///////////////////////////////////////////////////////////////////

package server

import (
	"fmt"

	"github.com/fasthttp/router"
	"github.com/futcity/controller/server/api"
	"github.com/futcity/controller/server/handlers"
	"github.com/valyala/fasthttp"
)

// WebServer Main server
type WebServer struct {
	relayh *handlers.RelayHandler
	grph   *handlers.GroupHandler
	devh   *handlers.DeviceHandler
	profh  *handlers.ProfileHandler
}

// NewWebServer Make new struct
func NewWebServer(rh *handlers.RelayHandler, gh *handlers.GroupHandler,
	dh *handlers.DeviceHandler, ph *handlers.ProfileHandler) *WebServer {
	return &WebServer{
		relayh: rh,
		grph:   gh,
		devh:   dh,
		profh:  ph,
	}
}

// IndexHandler Index page handler
func (w *WebServer) IndexHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetContentType("text/html")
	ctx.WriteString("<html><b>FutCity Controller</b></html>")
}

// NotFoundHandler Not found page handler
func (w *WebServer) NotFoundHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetContentType("text/html")
	ctx.WriteString("<html><b>Not found!</b></html>")
}

// Start Web server
func (w *WebServer) Start(ip string, port int) error {
	r := router.New()

	r.GET("/", w.IndexHandler)
	r.NotFound = w.NotFoundHandler

	r.GET(api.HttpReqGroupList, w.grph.Groups)
	r.GET(api.HttpReqDevList, w.devh.DeviceList)
	r.GET(api.HttpReqDevAdd, w.devh.AddDevice)
	r.GET(api.HttpReqDevRemove, w.devh.RemoveDevice)
	r.GET(api.HttpReqDevByDesc, w.devh.DeviceByDescription)

	r.GET(api.HttpReqProfAdd, w.profh.AddProfile)
	r.GET(api.HttpReqProfAddDev, w.profh.AddProfileDevice)
	r.GET(api.HttpReqProfAddGrp, w.profh.AddProfileGroup)
	r.GET(api.HttpReqProfDevRemove, w.profh.RemoveProfileDevice)
	r.GET(api.HttpReqProfGrpRemove, w.profh.RemoveProfileGroup)
	r.GET(api.HttpReqProfRemove, w.profh.RemoveProfile)
	r.GET(api.HttpReqProfList, w.profh.ProfileList)

	r.GET(api.HttpReqRelayStatus, w.relayh.Status)
	r.GET(api.HttpReqRelaySet, w.relayh.SetStatus)
	r.GET(api.HttpReqRelayUpdate, w.relayh.Update)
	r.GET(api.HttpReqRelaySwitch, w.relayh.Switch)
	r.GET(api.HttpReqRelayList, w.relayh.Devices)

	return fasthttp.ListenAndServe(fmt.Sprintf("%s:%d", ip, port), r.Handler)
}
