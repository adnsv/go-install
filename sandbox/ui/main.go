package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
)

var pp []*Paragraph

func render(rc *RC) {
	rc.Clear()
	sw, sh := rc.Size()

	w := 80
	if sw < 80 {
		w = sw
		if w < 32 {
			w = 32
		}
	}
	x := (sw - w) / 2
	if x < 0 {
		x = 0
	}

	h := 0
	for _, p := range pp {
		if p.WrapWidth != w {
			p.Wrap(w)
			h += p.Height + 1
		}
	}
	y := (sh - h + 1) / 2
	if y < 0 {
		y = 0
	}

	x = 0
	y = 0

	for _, p := range pp {
		p.Render(rc, x, y)
		y += p.Height + 1
	}
	//	rc.PutStr(w/2-7, h/2, tcell.StyleDefault, "Hello, World!")
	//	rc.PutStr(w/2-9, h/2+1, tcell.StyleDefault, "Press ESC to exit.")

	rc.Show()
}

func main() {
	encoding.Register()

	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlue).
		Foreground(tcell.ColorWhite)
	s.SetStyle(defStyle)

	pp := []*Paragraph{}
	pp = append(pp, NewParagraph("Hello, World!"))

	rc := &RC{s}
	render(rc)

	for {
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			rc.Sync()
			render(rc)
		case *tcell.EventKey:
			_, key, _ := ev.Modifiers(), ev.Key(), ev.Rune()
			switch key {
			case tcell.KeyEscape:
				{
					s.Fini()
					os.Exit(0)
				}
			case tcell.KeyUp:
				{

				}
			}
		}
	}
}
