package main

import (
	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
)

type RC struct {
	tcell.Screen
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
