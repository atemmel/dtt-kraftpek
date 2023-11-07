package main

import (
	"github.com/gdamore/tcell/v2"
	_ "github.com/gdamore/tcell/v2/encoding"
	"github.com/mattn/go-runewidth"
)

var (
	screen  tcell.Screen
	running bool = true
)

func Assert(err error) {
	if err != nil {
		panic(err)
	}
}

func DrawStr(screen tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		screen.SetContent(x, y, c, comb, style)
		x += w
	}
}

func Draw() {
	w, h := screen.Size()
	screen.Clear()
	style := tcell.StyleDefault.
		Foreground(tcell.ColorCadetBlue.TrueColor()).
		Background(tcell.ColorWhite)
	DrawStr(screen, w/2-7, h/2, style, "Hello, World!")
	DrawStr(screen, w/2-9, h/2+1, tcell.StyleDefault, "Press ESC to exit.")
	screen.Show()
}

func Quit() {
	screen.Fini()
	running = false
}

func main() {
	var err error
	screen, err = tcell.NewScreen()
	Assert(err)
	err = screen.Init()
	Assert(err)

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)

	Draw()

	for running {
		switch ev := screen.PollEvent().(type) {
		case *tcell.EventResize:
			screen.Sync()
			Draw()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				Quit()
			}
		}
	}
}
