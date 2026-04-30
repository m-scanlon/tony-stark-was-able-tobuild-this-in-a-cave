package reality

import "skyra-v05/src/debug"

type User struct {
	id        string
	Realities map[string]Reality
}

func (u *User) ID() string { return u.id }

func (u *User) Create(r *Relation) Reality {
	return &User{
		id:        r.ID,
		Realities: make(map[string]Reality),
	}
}

func (u *User) Realize(r *Relation) string {
	debug.Log("[user]: realizing", u.id)

	if being, ok := u.Realities["being"]; ok {
		being.Realize(r)
	}

	device, ok := u.Realities["device"]
	if !ok {
		return ""
	}

	result := device.Realize(r)
	debug.Log("[user]: device returned:", result)

	r.ID = r.Origin
	r.Origin = u.id
	return result
}

func (u *User) Parse() string {
	return ""
}
