///////////////////////////////////////////////////////////////////
//
// Future City Project
//
// Copyright (C) 2020-2021 Sergey Denisov. GPLv3
//
// Written by Sergey Denisov aka LittleBuster (DenisovS21@gmail.com)
//
///////////////////////////////////////////////////////////////////

package api

const (
	//
	// Common API
	//
	HttpReqGroupList = "/user/{user}/groups"

	//
	// Relay API
	//
	HttpReqRelayList   = "/user/{user}/relay"
	HttpReqRelayStatus = "/user/{user}/relay/{id}"
	HttpReqRelaySwitch = "/user/{user}/relay/{id}/switch"
	HttpReqRelaySet    = "/user/{user}/relay/{id}/set/{status}"
	HttpReqRelayUpdate = "/user/{user}/relay/{id}/update/state/{state}"

	//
	// Light API
	//
	HttpReqLightList     = "/user/{user}/light"
	HttpReqLightStatus   = "/user/{user}/light/{id}"
	HttpReqLightSwitch   = "/user/{user}/light/{id}/switch"
	HttpReqLightSet      = "/user/{user}/light/{id}/set/{status}"
	HttpReqLightGroupSet = "/user/{user}/light/{id}/group/{group}/set/{status}"
	HttpReqLightUpdate   = "/user/{user}/light/{id}/update/state/{state}"
)

//
// General responses
//

type DeviceResponse struct {
	Operation   string `json: "operation"`
	Result      bool   `json:"result"`
	Error       string `json:"error"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Online      bool   `json:"online"`
}

type GroupResponse struct {
	Operation string   `json: "operation"`
	Result    bool     `json:"result"`
	Error     string   `json:"error"`
	Groups    []string `json:"name"`
}

//
// Relay responses
//

type RelayResponse struct {
	Operation string `json: "operation"`
	Result    bool   `json: "result"`
	Error     string `json: "error"`
	Status    bool   `json:"status"`
	State     bool   `json:"state"`
}

type RelaySingleDevResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Online      bool   `json:"online"`
	Status      bool   `json:"status"`
	State       bool   `json:"state"`
}

type RelayDevResponse struct {
	Operation string                   `json: "operation"`
	Result    bool                     `json: "result"`
	Error     string                   `json: "error"`
	Relays    []RelaySingleDevResponse `json: "relays"`
}

//
// Light switch responses
//

type LightResponse struct {
	Operation string `json: "operation"`
	Result    bool   `json: "result"`
	Error     string `json: "error"`
	Status    bool   `json:"status"`
	State     bool   `json:"state"`
}

type LightSingleDevResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Online      bool   `json:"online"`
	Status      bool   `json:"status"`
	State       bool   `json:"state"`
}

type LightDevResponse struct {
	Operation string                   `json: "operation"`
	Result    bool                     `json: "result"`
	Error     string                   `json: "error"`
	Lights    []LightSingleDevResponse `json: "lights"`
}
