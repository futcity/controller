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
	"github.com/futcity/controller/server/api"
	"github.com/futcity/controller/utils"
	"github.com/valyala/fasthttp"
)

type GroupHandler struct {
	aut *auth.Authorization
	log *utils.Log
}

func NewGroupHandler(a *auth.Authorization, l *utils.Log) *GroupHandler {
	return &GroupHandler{
		aut: a,
		log: l,
	}
}

func (r *GroupHandler) Groups(ctx *fasthttp.RequestCtx) {
	// Check user rights
	var groups, err = r.aut.Groups(ctx.UserValue("user").(string))
	if err != nil {
		r.responseList(ctx, "Groups list", false, "Authorization failed", nil)
		return
	}

	// Send response
	r.responseList(ctx, "Groups list", true, "", groups)
}

func (r *GroupHandler) responseList(ctx *fasthttp.RequestCtx, oper string, result bool, err string, groups *[]string) {
	ctx.Response.Header.SetContentType("application/json")

	var grpResp = api.GroupResponse{
		Operation: oper,
		Result:    result,
		Error:     err,
	}

	for _, group := range *groups {
		grpResp.Groups = append(grpResp.Groups, group)
	}

	if result {
		r.log.Info("GROUPH", oper)
	} else {
		r.log.Error("GROUPH", oper, err)
	}

	var bytes, _ = json.Marshal(grpResp)

	ctx.Write(bytes)
}
