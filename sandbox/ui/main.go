package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
)

var pp []*Paragraph

func render(rc *RC) {
	rc.Clear()
	sw, sh := rc.Size()

	w := sw - 10
	if w > 80 {
		w = 80
	} else if w < 32 {
		w = 32
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

	for _, p := range strings.Split(license, "\n") {
		pp = append(pp, NewParagraph(p))
	}

	rc := &RC{s, false}
	render(rc)

	for {
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			rc.ScheduleUpdate()

		case *EventUpdate:
			rc.ClearUpdatePending()
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

var copyright = "Copyright 2021 Tessonics Inc. All Rights Reserved"
var license = `Please read this software license and service agreement ("License") carefully before accepting the license conditions. By accepting you ("End-User") are agreeing to be bound by the terms of this license. If you do not agree to the terms of this license, please cancel the installation process.
The software is provided by Tessonics Inc. ("Tessonics")
1. License. Tessonics grants the End-User a limited, non-exclusive, nontransferable, royalty-free license to install and use the included software, documentation and any other files accompanying this License whether on disk, in read only memory, on any other media or in any other form.
2. Permitted Uses and Restrictions. This License allows the End-User to install and use the Software. Except as expressly permitted in this License, the End-User may not decompile, reverse engineer, disassemble, modify, rent, lease, loan, sublicense, distribute or create derivative works based upon the Software in whole or in part or transmit the Software over a network. The software is not intended for use in any operation of any facilities, in which case the failure of the software could lead to death, personal injury, or severe physical or environmental damage. The End-User's rights under this License will terminate automatically without notice if the End-User fails to comply with any term(s) of this License. This License shall terminate upon the End User's cessation to use the Software. Upon termination of this License, the End User shall uninstall all the Tessonics software.
3. To the maximum extent permitted by applicable law, in no event shall Tessonics or its suppliers be liable for any consequential, incidental, direct, indirect, special, punitive, or other damages whatsoever (including, without limitation, damages for loss of business profits, business interruption, loss of business information, or other pecuniary loss) arising out of the use of or inability to use the Software or documentation, even if Tessonics has been advised of the possibility of such damages. Some jurisdictions do not allow the exclusion or limitation of liability for consequential or incidental damages, so the above limitation may not apply to the end-user. 
4. Disclaimer of Warranty on Software. The End-User expressly acknowledges and agrees that use of the Software is at the End-User's sole risk. The Software is provided "AS IS" and without warranty of any kind and Tessonics expressly disclaims all warranties and/or conditions, express or implied, including, but not limited to, the implied warranties and/or conditions of merchantability or satisfactory quality and fitness for a particular purpose and noninfringement of third party rights. Tessonics does not warrant that the functions contained in the Software will meet the end-user's requirements, or that the operation of the Software will be uninterrupted or error-free, or that defects in the Software will be corrected. furthermore, Tessonics does not warrant or make any representations regarding the use or the results of the use of the Software or related documentation in terms of their correctness, accuracy, reliability, or otherwise. No oral or written information or advice given by Tessonics or any agent shall create a warranty or in any way increase the scope of this warranty. Should the Software prove defective, the end-user assumes the entire cost of all necessary servicing, repair or correction. Some jurisdictions do not allow the exclusion of implied warranties, so the above exclusion may not apply to the end-user. 
5. Hacking etc. Tessonics shall not be responsible or liable to Customer or any End User for any acts of fraud, theft, misappropriation, tampering, hacking, interception, piracy, misuse, misrepresentation, dissemination, or other legal or unauthorized activities of third parties.
6. Complete Agreement. This License constitutes the entire agreement between the parties with respect to the use of the Software and supersedes all prior or contemporaneous understandings regarding such subject matter. No amendment to or modification of this License will be binding unless in writing and signed by Tessonics.
`
