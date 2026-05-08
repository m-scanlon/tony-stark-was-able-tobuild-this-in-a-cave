package reality

import (
	"fmt"
	"skyra-v05/src/debug"
)

const xpPerEntry = 5

var levelThresholds = []int{0, 50, 150, 300, 500, 750, 1050, 1400, 1800, 2250}

type Levels struct {
	id string
	XP map[string]int
}

func (l *Levels) ID() string { return l.id }

func (l *Levels) Create(r *Relation) Reality {
	return &Levels{
		id: "levels",
		XP: make(map[string]int),
	}
}

func (l *Levels) Award(sender, recipient string) {
	l.XP[sender] += xpPerEntry
	l.XP[recipient] += xpPerEntry
	debug.Log("[levels]: +", xpPerEntry, "xp →", sender, "(", l.XP[sender], ") +", recipient, "(", l.XP[recipient], ")")
}

func (l *Levels) Level(name string) int {
	xp := l.XP[name]
	level := 1
	for i, threshold := range levelThresholds {
		if xp >= threshold {
			level = i + 1
		}
	}
	return level
}

func (l *Levels) Realize(r *Relation) string {
	if r.Collecting {
		snap := make(map[string]LevelSnapshot)
		for name, xp := range l.XP {
			snap[name] = LevelSnapshot{
				XP:    xp,
				Level: l.Level(name),
				Next:  l.xpToNext(name),
			}
		}
		r.Export("levels", snap)
		return ""
	}
	return ""
}

func (l *Levels) ParseFor(self, peer string) string {
	return fmt.Sprintf("your level: %d (xp: %d, next: %d)\n", l.Level(self), l.XP[self], l.xpToNext(self))
}

func (l *Levels) xpToNext(name string) int {
	xp := l.XP[name]
	for _, threshold := range levelThresholds {
		if xp < threshold {
			return threshold - xp
		}
	}
	return 0
}
