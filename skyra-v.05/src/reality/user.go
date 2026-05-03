package reality

import "skyra-v05/src/debug"

type User struct {
	id        string
	Realities map[string]Reality
}

func (u *User) ID() string { return u.id }

func (u *User) Create(r *Relation) Reality {
	user := &User{
		id:        r.ID,
		Realities: make(map[string]Reality),
	}

	if r.Impulse != "" {
		being := Being{}.Create(r).(Being)
		user.Realities["being"] = being
	}

	if r.Realities != nil {
		for _, reality := range r.Realities {
			if _, ok := reality.(*MacOS); ok {
				user.Realities["device"] = reality
				break
			}
		}
	}

	return user
}

func (u *User) Realize(r *Relation) string {
	if r.Collecting {
		node := RealityNode{ID: u.id, Type: "User", Children: []RealityNode{}}
		snap := BeingSnapshot{
			Name: u.id, Type: "user", Status: "idle",
			Peers:    []string{},
			Memories: MemorySnapshot{Items: []MemoryItem{}, Skills: []SkillItem{}},
		}

		if being, ok := u.Realities["being"].(Being); ok {
			snap.Identity = being.Identity
			snap.Purpose = being.Purpose
			if being.Relationships != nil {
				snap.Peers = being.Relationships
			}
			snap.Device = being.Device
			snap.Memories = snapshotMemories(being.Home)
			node.Children = append(node.Children, RealityNode{ID: u.id + "-being", Type: "Being", Children: []RealityNode{}})
		}

		if device, ok := u.Realities["device"]; ok {
			node.Children = append(node.Children, RealityNode{
				ID: device.ID(), Type: capitalizeType(device.ID()), Children: []RealityNode{},
			})
		}

		r.Export("being:"+u.id, snap)
		r.Export("node:"+u.id, node)
		return ""
	}

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

	r.Impulse = result
	r.ID = r.Origin
	r.Origin = u.id
	return result
}
