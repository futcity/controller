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

import "github.com/futcity/controller/core/devices"

// Storage All devices map
type Storage struct {
	devices map[string]devices.IDevice
}

// NewStorage Make new storage
func NewStorage() *Storage {
	return &Storage{
		devices: make(map[string]devices.IDevice),
	}
}

// AddDevice Add new device in storage
func (s *Storage) AddDevice(dev devices.IDevice) {
	s.devices[dev.Name()] = dev
}

// Device Get device by name
func (s *Storage) Device(name string) devices.IDevice {
	return s.devices[name]
}

// Devices Get all devices
func (s *Storage) Devices() []devices.IDevice {
	var list []devices.IDevice

	for _, dev := range s.devices {
		list = append(list, dev)
	}

	return list
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
