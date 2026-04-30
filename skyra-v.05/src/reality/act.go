package reality

type Act struct {
	id        string
	Realities map[string]Reality
}

func (a *Act) ID() string { return a.id }

func (a *Act) Create(r *Relation) Reality {
	return &Act{
		id:        "act",
		Realities: make(map[string]Reality),
	}
}

func (a *Act) Realize(r *Relation) string {
	r.Attach("act", a.Parse)
	if device, ok := a.Realities["device"]; ok {
		return device.Realize(r)
	}
	return ""
}

func (a *Act) Parse() string {
	return ""
}
