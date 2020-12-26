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

type ProfDevCfg struct {
	Name  string
	Read  bool
	Write bool
}

type ProfileCfg struct {
	Name    string
	Key     string
	Admin   bool
	Groups  []string
	Devices []ProfDevCfg
}

type ProfCfg struct {
	Profiles []ProfileCfg
}
