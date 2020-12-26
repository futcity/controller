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

type DbFileCfg struct {
	Name string
	Path string
}

type DbCfg struct {
	Type  string
	Files []DbFileCfg
}

type ServerCfg struct {
	IP   string
	Port int
}

type AppCfg struct {
	Server ServerCfg
	Db     DbCfg
}
