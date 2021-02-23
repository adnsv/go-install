package main

import (
	"time"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
)

type RC struct {
	tcell.Screen
	updatePending bool
}

func (rc *RC) PutStr(x, y int, style tcell.Style, s string) {
	rc.PutRunes(x, y, style, []rune(s))
}

func (rc *RC) PutRunes(x, y int, style tcell.Style, runes []rune) {
	for _, r := range runes {
		var comb []rune
		w := runewidth.RuneWidth(r)
		if w == 0 {
			comb = []rune{r}
			r = ' '
			w = 1
		}
		rc.SetContent(x, y, r, comb, style)
		x += w
	}
}

func (rc *RC) ScheduleUpdate() {
	if !rc.updatePending {
		rc.updatePending = true
		rc.PostEvent(&EventUpdate{t: time.Now()})
	}
}

func (rc *RC) ClearUpdatePending() {
	rc.updatePending = false
}

type EventUpdate struct {
	t time.Time
}

func (ev *EventUpdate) When() time.Time {
	return ev.t
}
