///////////////////////////////////////////////////////////////////
//
// Future City Project
//
// Copyright (C) 2020-2021 Sergey Denisov. GPLv3
//
// Written by Sergey Denisov aka LittleBuster (DenisovS21@gmail.com)
//
///////////////////////////////////////////////////////////////////

package core

import (
	"errors"

	"github.com/futcity/controller/core/devices"
	"github.com/futcity/controller/core/devices/base"
	"github.com/futcity/controller/utils"
)

// Storage All devices map
type Storage struct {
	devices map[string]devices.IDevice
	log     *utils.Log
}

// NewStorage Make new storage
func NewStorage(log *utils.Log) *Storage {
	return &Storage{
		devices: make(map[string]devices.IDevice),
		log:     log,
	}
}

// AddDevice Add new device in storage
func (s *Storage) AddDevice(name string, desc string, devType string) {
	var device devices.IDevice

	if devType == "relay" {
		device = base.NewRelay(name, desc)
	}
	device.SetID(len(s.devices))
	s.devices[name] = device
}

func (s *Storage) RemoveByID(id int) error {
	for _, dev := range s.devices {
		if dev.ID() == id {
			delete(s.devices, dev.Name())
			return nil
		}
	}
	return errors.New("Device not found")
}

// Device Get device by name
func (s *Storage) Device(name string) devices.IDevice {
	return s.devices[name]
}

// DevicesByType Get all devices by type
func (s *Storage) DevicesByType(devType string) []devices.IDevice {
	var list []devices.IDevice

	for _, dev := range s.devices {
		if dev.Type() == devType {
			list = append(list, dev)
		}
	}

	return list
}

// DeviceByDescription Get device by description
func (s *Storage) DeviceByDescription(desc string) (devices.IDevice, error) {
	for _, dev := range s.devices {
		if dev.Description() == desc {
			return dev, nil
		}
	}
	return nil, errors.New("Device not found")
}

func (s *Storage) DeviceByID(id int) (devices.IDevice, error) {
	for _, dev := range s.devices {
		if dev.ID() == id {
			return dev, nil
		}
	}
	return nil, errors.New("Device not found")
}

// Devices Get all devices
func (s *Storage) Devices() []devices.IDevice {
	var list []devices.IDevice

	for _, dev := range s.devices {
		list = append(list, dev)
	}

	return list
}
