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

import "errors"

// Authorization User profiles
type Authorization struct {
	prof map[string]*Profile
}

// NewAuthorization Make new struct
func NewAuthorization() *Authorization {
	return &Authorization{
		prof: make(map[string]*Profile),
	}
}

// AddProfile Add new user profile
func (a *Authorization) AddProfile(prof *Profile) {
	a.prof[prof.APIKey()] = prof
}

// DeleteProfile Delete user profile
func (a *Authorization) DeleteProfile(name string) error {
	for _, prof := range a.prof {
		if prof.Name() == name {
			delete(a.prof, name)
			return nil
		}
	}
	return errors.New("Profile not found")
}

// Validation Check user device by key
func (a *Authorization) Validation(key string, device string) (bool, bool) {
	var prof = a.prof[key]
	if prof == nil {
		return false, false
	}

	if prof.Admin() {
		return true, true
	}

	var dev = prof.Device(device)
	if dev == nil {
		return false, false
	}

	return dev.Read(), dev.Write()
}

// Profiles Get all profiles
func (a *Authorization) Profiles() []*Profile {
	var profiles []*Profile

	for _, prof := range a.prof {
		profiles = append(profiles, prof)
	}

	return profiles
}

// Profiles Get profile by name
func (a *Authorization) Profile(name string) *Profile {
	for _, profile := range a.prof {
		if profile.Name() == name {
			return profile
		}
	}
	return nil
}

// IsAdmin Check admin status
func (a *Authorization) IsAdmin(key string) bool {
	var prof = a.prof[key]
	if prof == nil {
		return false
	}

	if prof.Admin() {
		return true
	}

	return false
}

// Groups Get profile groups
func (a *Authorization) Groups(key string) (*[]string, error) {
	var profile = a.prof[key]
	if profile == nil {
		return nil, errors.New("Profile not found")
	}

	return profile.Groups(), nil
}
