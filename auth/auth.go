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

// Groups Get profile groups
func (a *Authorization) Groups(key string) (*[]string, error) {
	var profile = a.prof[key]
	if profile == nil {
		return nil, errors.New("Profile not found")
	}

	return profile.Groups(), nil
}
