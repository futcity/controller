///////////////////////////////////////////////////////////////////
//
// Future City Project
//
// Copyright (C) 2020-2021 Sergey Denisov. GPLv3
//
// Written by Sergey Denisov aka LittleBuster (DenisovS21@gmail.com)
//
///////////////////////////////////////////////////////////////////

package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Configs App configs
type Configs struct {
}

// NewConfigs make new struct
func NewConfigs() *Configs {
	return &Configs{}
}

// LoadFromFile Loading configs from file
func (c *Configs) LoadFromFile(settings interface{}, fileName string) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, settings)
	if err != nil {
		return err
	}

	return nil
}

// SaveToFile Save configs to file
func (c *Configs) SaveToFile(settings interface{}, fileName string) error {
	var bytes, err = json.Marshal(settings)
	if err != nil {
		return err
	}

	var file, err2 = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0660)
	if err2 != nil {
		return err2
	}
	defer file.Close()

	_, err = file.WriteString(string(bytes))
	if err != nil {
		return err
	}

	return nil
}
