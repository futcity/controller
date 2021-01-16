package db

type ProfileDeivceDB struct {
	Name  string `json:"name"`
	Read  bool   `json:"read"`
	Write bool   `json:"write"`
}

type SingleProfileDB struct {
	Name    string            `json:"name"`
	Key     string            `json:"key"`
	Admin   bool              `json:"admin"`
	Groups  []string          `json:"groups"`
	Devices []ProfileDeivceDB `json:"devices"`
}

type ProfileDB struct {
	Profiles []SingleProfileDB `json:"profiles"`
}
