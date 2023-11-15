package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
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

func TuiPanic(screen tcell.Screen) {
	if r := recover(); r != nil {
		screen.Fini()
		panic(r)
	}
}
