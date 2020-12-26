///////////////////////////////////////////////////////////////////
//
// Future City Project
//
// Copyright (C) 2020-2021 Sergey Denisov. GPLv3
//
// Written by Sergey Denisov aka LittleBuster (DenisovS21@gmail.com)
//
///////////////////////////////////////////////////////////////////

package configs

type DeviceCfg struct {
	Name        string
	Description string
	Type        string
}

type DevCfg struct {
	Devices []DeviceCfg
}
