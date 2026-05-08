package reality

import "strings"

type Plan struct {
	id string
}

func (p *Plan) ID() string { return p.id }

func (p *Plan) Create(r *Relation) Reality {
	return &Plan{id: "plan"}
}

func (p *Plan) Realize(r *Relation) string {
	log := r.Log

	desk := findDesk(r)
	if desk == nil {
		if log != nil {
			log("[plan]: no desk on relation")
		}
		return "no desk available"
	}

	impulse := strings.TrimSpace(r.Impulse)
	if impulse == "" {
		return "no plan command"
	}

	if cmd, err := ExtractTag(impulse, "create-task"); err == nil {
		return p.createTask(desk, cmd, log)
	}
	if cmd, err := ExtractTag(impulse, "complete-task"); err == nil {
		return p.completeTask(desk, cmd, log)
	}
	if cmd, err := ExtractTag(impulse, "drop-task"); err == nil {
		return p.dropTask(desk, cmd, log)
	}
	if cmd, err := ExtractTag(impulse, "open-task"); err == nil {
		return p.openTask(desk, cmd, log)
	}
	if cmd, err := ExtractTag(impulse, "close-task"); err == nil {
		return p.closeTask(desk, cmd, log)
	}
	if cmd, err := ExtractTag(impulse, "focus-task"); err == nil {
		return p.focusTask(desk, cmd, log)
	}

	return "unknown plan command"
}

func (p *Plan) createTask(desk *Desk, cmd string, log func(...any)) string {
	rel, err := ExtractTag(cmd, "relationship")
	if err != nil {
		return "create-task: missing <relationship>"
	}
	name, err := ExtractTag(cmd, "name")
	if err != nil {
		return "create-task: missing <name>"
	}

	task := &Task{Name: name, State: "active"}

	if desc, err := ExtractTag(cmd, "description"); err == nil {
		task.Description = desc
	}
	if val, err := ExtractTag(cmd, "validation"); err == nil {
		task.Validation = val
	}
	if assumptions, err := ExtractTag(cmd, "assumptions"); err == nil {
		for _, a := range strings.Split(assumptions, ",") {
			a = strings.TrimSpace(a)
			if a != "" {
				task.Assumptions = append(task.Assumptions, a)
			}
		}
	}
	if commands, err := ExtractTag(cmd, "commands"); err == nil {
		for _, c := range strings.Split(commands, ",") {
			c = strings.TrimSpace(c)
			if c != "" {
				task.Commands = append(task.Commands, c)
			}
		}
	}

	if parent, err := ExtractTag(cmd, "parent"); err == nil {
		parentTask := desk.FindTask(rel, parent)
		if parentTask == nil {
			return "create-task: parent " + parent + " not found in " + rel
		}
		parentTask.Items = append(parentTask.Items, task)
		desk.OpenTask(rel, parent)
	} else {
		desk.CreateTask(rel, task)
	}

	if log != nil {
		log("[plan]: created task", name, "in", rel)
	}
	return "created: " + name + " [" + rel + "]"
}

func (p *Plan) completeTask(desk *Desk, cmd string, log func(...any)) string {
	rel, err := ExtractTag(cmd, "relationship")
	if err != nil {
		return "complete-task: missing <relationship>"
	}
	name, err := ExtractTag(cmd, "name")
	if err != nil {
		return "complete-task: missing <name>"
	}
	if err := desk.CompleteTask(rel, name); err != nil {
		return "complete-task: " + err.Error()
	}
	if log != nil {
		log("[plan]: submitted for review", name, "in", rel)
	}
	return "submitted for review: " + name + " [" + rel + "]"
}

func (p *Plan) dropTask(desk *Desk, cmd string, log func(...any)) string {
	rel, err := ExtractTag(cmd, "relationship")
	if err != nil {
		return "drop-task: missing <relationship>"
	}
	name, err := ExtractTag(cmd, "name")
	if err != nil {
		return "drop-task: missing <name>"
	}
	if err := desk.DropTask(rel, name); err != nil {
		return "drop-task: " + err.Error()
	}
	if log != nil {
		log("[plan]: dropped task", name, "in", rel)
	}
	return "dropped: " + name + " [" + rel + "]"
}

func (p *Plan) openTask(desk *Desk, cmd string, log func(...any)) string {
	rel, err := ExtractTag(cmd, "relationship")
	if err != nil {
		return "open-task: missing <relationship>"
	}
	name, err := ExtractTag(cmd, "name")
	if err != nil {
		return "open-task: missing <name>"
	}
	desk.OpenTask(rel, name)
	if log != nil {
		log("[plan]: opened task", name, "in", rel)
	}
	return "opened: " + name + " [" + rel + "]"
}

func (p *Plan) closeTask(desk *Desk, cmd string, log func(...any)) string {
	rel, err := ExtractTag(cmd, "relationship")
	if err != nil {
		return "close-task: missing <relationship>"
	}
	name, err := ExtractTag(cmd, "name")
	if err != nil {
		return "close-task: missing <name>"
	}
	desk.CloseTask(rel, name)
	if log != nil {
		log("[plan]: closed task", name, "in", rel)
	}
	return "closed: " + name + " [" + rel + "]"
}

func (p *Plan) focusTask(desk *Desk, cmd string, log func(...any)) string {
	rel, err := ExtractTag(cmd, "relationship")
	if err != nil {
		return "focus-task: missing <relationship>"
	}
	name, err := ExtractTag(cmd, "name")
	if err != nil {
		return "focus-task: missing <name>"
	}
	desk.FocusTask(rel, name)
	if log != nil {
		log("[plan]: focused task", name, "in", rel)
	}
	return "focused: " + name + " [" + rel + "]"
}

func findDesk(r *Relation) *Desk {
	if r.Realities == nil {
		return nil
	}
	if d, ok := r.Realities["desk"]; ok {
		if desk, ok := d.(*Desk); ok {
			return desk
		}
	}
	return nil
}
