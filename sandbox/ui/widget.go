package main

import (
	"strings"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
)

type Paragraph struct {
	WrapWidth int
	Height    int
	Words     []*Word
	Lines     [][]*Word
}

type Word struct {
	Content      []rune
	ContentWidth int
}

func NewParagraph(text string) *Paragraph {
	p := &Paragraph{}
	for _, s := range strings.Split(text, " ") {
		if s == "" {
			continue
		}
		rr := []rune(s)
		w := 0
		for _, r := range rr {
			w += runewidth.RuneWidth(r)
		}
		p.Words = append(p.Words, &Word{
			Content:      rr,
			ContentWidth: w,
		})
	}
	return p
}

func (p *Paragraph) Wrap(maxw int) {
	p.WrapWidth = maxw
	p.Lines = p.Lines[:0]
	if len(p.Words) == 0 {
		return
	}
	ww := p.Words
	line := []*Word{ww[0]}
	lw := ww[0].ContentWidth
	ww = ww[1:]
	for len(ww) > 0 {
		word := ww[0]
		ww = ww[1:]
		nextw := word.ContentWidth
		if lw+1+nextw > maxw {
			p.Lines = append(p.Lines, line)
			line = []*Word{word}
			lw = nextw
		} else {
			line = append(line, word)
			lw += 1 + nextw
		}
	}
	p.Lines = append(p.Lines, line)
	p.Height = len(p.Lines)
}

func (p *Paragraph) Render(rc *RC, x, y int) {
	for _, ln := range p.Lines {
		cx := x
		for _, w := range ln {
			rc.PutRunes(cx, y, tcell.StyleDefault, w.Content)
			cx += w.ContentWidth + 1
		}
		y++
	}
}
