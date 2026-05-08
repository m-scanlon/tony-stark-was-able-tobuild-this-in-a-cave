package reality

import (
	"fmt"
	"skyra-v05/src/debug"
	"strings"
)


type Desk struct {
	id                string
	Owner             string
	Items             map[string][]*Task
	Views             map[string]string
	RelationshipViews map[string]string
}

func (d *Desk) ID() string { return d.id }

func (d *Desk) Create(r *Relation) Reality {
	return &Desk{
		id:                "desk",
		Items:             make(map[string][]*Task),
		Views:             make(map[string]string),
		RelationshipViews: make(map[string]string),
	}
}

func (d *Desk) Realize(r *Relation) string {
	if r.Collecting {
		snap := d.Snapshot()
		r.Export("desk:"+d.Owner, snap)
		return ""
	}
	if d.Empty() {
		return ""
	}
	r.Attach("desk", d.Parse)
	return ""
}

func (d *Desk) Empty() bool {
	return len(d.Items) == 0
}

func (d *Desk) Parse() string {
	if len(d.Items) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("desk:\n")
	for rel, tasks := range d.Items {
		view := d.RelationshipViews[rel]
		if view == "collapsed" {
			sb.WriteString("  " + rel + ": (collapsed)\n")
			continue
		}
		sb.WriteString("  " + rel + ":\n")
		for _, task := range tasks {
			open := d.Views[viewKey(rel, task.Name)] == "open"
			sb.WriteString(task.Realize(4, open))
		}
	}
	return sb.String()
}

func (d *Desk) CreateTask(relationship string, task *Task) {
	if task.State == "" {
		task.State = "active"
	}
	d.Items[relationship] = append(d.Items[relationship], task)
	if _, ok := d.RelationshipViews[relationship]; !ok {
		d.RelationshipViews[relationship] = "open"
	}
	debug.Being(d.Owner, "desk", "[desk]: created task", task.Name, "in", relationship, "state:", task.State)
}

func (d *Desk) FindTask(relationship, name string) *Task {
	for _, task := range d.Items[relationship] {
		if found := findInTree(task, name); found != nil {
			return found
		}
	}
	return nil
}

func findInTree(task *Task, name string) *Task {
	if strings.EqualFold(task.Name, name) {
		return task
	}
	for _, sub := range task.Items {
		if found := findInTree(sub, name); found != nil {
			return found
		}
	}
	return nil
}

func (d *Desk) CompleteTask(relationship, name string) error {
	task := d.FindTask(relationship, name)
	if task == nil {
		return fmt.Errorf("task %q not found in %s", name, relationship)
	}
	task.State = "review"
	debug.Being(d.Owner, "desk", "[desk]: completed", name, "in", relationship, "→ review")
	return nil
}

func (d *Desk) AcceptTask(relationship, name, acceptor string) error {
	task := d.FindTask(relationship, name)
	if task == nil {
		return fmt.Errorf("task %q not found in %s", name, relationship)
	}
	if task.State != "review" {
		return fmt.Errorf("task %q is not in review", name)
	}
	task.State = "done"
	task.AcceptedBy = acceptor
	debug.Being(d.Owner, "desk", "[desk]: accepted", name, "by", acceptor)
	return nil
}

func (d *Desk) RejectTask(relationship, name string) error {
	task := d.FindTask(relationship, name)
	if task == nil {
		return fmt.Errorf("task %q not found in %s", name, relationship)
	}
	if task.State != "review" {
		return fmt.Errorf("task %q is not in review", name)
	}
	task.State = "active"
	debug.Being(d.Owner, "desk", "[desk]: rejected", name, "→ active")
	return nil
}

func (d *Desk) DropTask(relationship, name string) error {
	task := d.FindTask(relationship, name)
	if task == nil {
		return fmt.Errorf("task %q not found in %s", name, relationship)
	}
	task.State = "dropped"
	debug.Being(d.Owner, "desk", "[desk]: dropped", name)
	return nil
}

func (d *Desk) OpenTask(relationship, name string) {
	d.Views[viewKey(relationship, name)] = "open"
}

func (d *Desk) CloseTask(relationship, name string) {
	d.Views[viewKey(relationship, name)] = "closed"
}

func (d *Desk) FocusTask(relationship, name string) {
	for _, task := range d.Items[relationship] {
		d.Views[viewKey(relationship, task.Name)] = "closed"
	}
	d.Views[viewKey(relationship, name)] = "open"
}

func (d *Desk) OpenRelationship(relationship string) {
	d.RelationshipViews[relationship] = "open"
}

func (d *Desk) CollapseRelationship(relationship string) {
	d.RelationshipViews[relationship] = "collapsed"
}

func (d *Desk) ParseScoped(relationships []string) string {
	if len(d.Items) == 0 {
		return ""
	}
	scope := make(map[string]bool)
	for _, r := range relationships {
		scope[r] = true
	}
	var sb strings.Builder
	sb.WriteString("desk:\n")
	for rel, tasks := range d.Items {
		if !scope[rel] {
			continue
		}
		sb.WriteString("  " + rel + ":\n")
		for _, task := range tasks {
			open := d.Views[viewKey(rel, task.Name)] == "open"
			sb.WriteString(task.Realize(4, open))
		}
	}
	return sb.String()
}

func viewKey(relationship, name string) string {
	return relationship + ":" + name
}

func (d *Desk) Snapshot() DeskSnapshot {
	snap := DeskSnapshot{
		Tasks: make(map[string][]TaskSnapshot),
		Views: make(map[string]string),
	}
	for rel, tasks := range d.Items {
		for _, task := range tasks {
			snap.Tasks[rel] = append(snap.Tasks[rel], snapshotTask(task))
		}
	}
	for k, v := range d.Views {
		snap.Views[k] = v
	}
	for k, v := range d.RelationshipViews {
		snap.Views["rel:"+k] = v
	}
	return snap
}

func snapshotTask(t *Task) TaskSnapshot {
	snap := TaskSnapshot{
		Name:        t.Name,
		Description: t.Description,
		Validation:  t.Validation,
		AcceptedBy:  t.AcceptedBy,
		State:       t.State,
	}
	if len(t.Assumptions) > 0 {
		snap.Assumptions = t.Assumptions
	}
	if len(t.Commands) > 0 {
		snap.Commands = t.Commands
	}
	for _, sub := range t.Items {
		snap.Items = append(snap.Items, snapshotTask(sub))
	}
	return snap
}
