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
	HttpReqDevList   = "/user/{user}/device"
	HttpReqDevRemove = "/user/{user}/device/del/id/{id}"
	HttpReqDevAdd    = "/user/{user}/device/add/name/{name}/desc/{desc}/type/{type}"
	HttpReqDevByDesc = "/user/{user}/device/desc/{desc}"

	HttpReqProfList      = "/user/{user}/profile"
	HttpReqProfAdd       = "/user/{user}/profile/add/name/{name}/key/{key}/admin/{admin}"
	HttpReqProfRemove    = "/user/{user}/profile/del/name/{name}"
	HttpReqProfAddDev    = "/user/{user}/profile/name/{name}/add/device/{device}/read/{read}/write/{write}"
	HttpReqProfAddGrp    = "/user/{user}/profile/name/{name}/add/group/{group}"
	HttpReqProfDevRemove = "/user/{user}/profile/name/{name}/del/device/{device}"
	HttpReqProfGrpRemove = "/user/{user}/profile/name/{name}/del/group/{group}"

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

// Devices responses

type DeviceResponse struct {
	Operation   string `json:"operation"`
	Result      bool   `json:"result"`
	Error       string `json:"error"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Online      bool   `json:"online"`
}

type DeviceSingleResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Online      bool   `json:"online"`
}

type DeviceListResponse struct {
	Operation string                 `json:"operation"`
	Result    bool                   `json:"result"`
	Error     string                 `json:"error"`
	Devices   []DeviceSingleResponse `json:"devices"`
}

// Groups response

type GroupResponse struct {
	Operation string   `json:"operation"`
	Result    bool     `json:"result"`
	Error     string   `json:"error"`
	Groups    []string `json:"groups"`
}

// Profiles responses

type ProfileDeviceResponse struct {
	Name  string `json:"name"`
	Read  bool   `json:"read"`
	Write bool   `json:"write"`
}

type ProfileSingleResponse struct {
	Name    string                  `json:"name"`
	Key     string                  `json:"key"`
	Admin   bool                    `json:"admin"`
	Groups  []string                `json:"groups"`
	Devices []ProfileDeviceResponse `json:"devices"`
}

type ProfileListResponse struct {
	Operation string                  `json:"operation"`
	Result    bool                    `json:"result"`
	Error     string                  `json:"error"`
	Profiles  []ProfileSingleResponse `json:"profiles"`
}

//
// Relay responses
//

type RelayResponse struct {
	Operation string `json:"operation"`
	Result    bool   `json:"result"`
	Error     string `json:"error"`
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
	Operation string                   `json:"operation"`
	Result    bool                     `json:"result"`
	Error     string                   `json:"error"`
	Relays    []RelaySingleDevResponse `json:"relays"`
}

//
// Light switch responses
//

type LightResponse struct {
	Operation string `json:"operation"`
	Result    bool   `json:"result"`
	Error     string `json:"error"`
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
	Operation string                   `json:"operation"`
	Result    bool                     `json:"result"`
	Error     string                   `json:"error"`
	Lights    []LightSingleDevResponse `json:"lights"`
}
