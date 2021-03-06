///////////////////////////////////////////////////////////////////
//
// Future City Project
//
// Copyright (C) 2020-2021 Sergey Denisov. GPLv3
//
// Written by Sergey Denisov aka LittleBuster (DenisovS21@gmail.com)
//
///////////////////////////////////////////////////////////////////

package auth

type Profile struct {
	name    string
	admin   bool
	key     string
	groups  []string
	devices map[string]*ProfileDevice
}

func NewProfile(name string, key string, admin bool) *Profile {
	return &Profile{
		name:    name,
		admin:   admin,
		key:     key,
		devices: make(map[string]*ProfileDevice),
	}
}

func (p *Profile) AddDevice(device *ProfileDevice) {
	p.devices[device.Name()] = device
}

func (p *Profile) AddGroup(grp string) {
	p.groups = append(p.groups, grp)
}

func (p *Profile) Groups() *[]string {
	return &p.groups
}

func (p *Profile) Device(name string) *ProfileDevice {
	return p.devices[name]
}

func (p *Profile) Devices() []*ProfileDevice {
	var devs []*ProfileDevice

	for _, dev := range p.devices {
		devs = append(devs, dev)
	}

	return devs
}

func (p *Profile) RemoveDevice(name string) {
	delete(p.devices, name)
}

func (p *Profile) RemoveGroup(name string) {
	var grps []string

	for _, grp := range p.groups {
		if grp != name {
			grps = append(grps, grp)
		}
	}
	p.groups = make([]string, len(grps))
	copy(p.groups, grps)
}

func (p *Profile) Name() string {
	return p.name
}

func (p *Profile) Admin() bool {
	return p.admin
}

func (p *Profile) APIKey() string {
	return p.key
}
