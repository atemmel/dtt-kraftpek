package main

import (
	"os"

	"github.com/gdamore/tcell/v2"
	_ "github.com/gdamore/tcell/v2/encoding"
	"github.com/mattn/go-runewidth"
)

var (
	screen tcell.Screen

	currentSlide = 0
	slides       = [][]string{
		{
			"Rubrik till slide 1",
			"",
			"* Punkt 1",
			"* Punkt 2",
			"* Punkt 3",
		},
		{
			"Slide 2",
		},
		{
			"Slide 3",
		},
	}
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
	screen.Clear()
	// exempel på hur styling kan sättas

	//style := tcell.StyleDefault.
	//Foreground(tcell.ColorCadetBlue.TrueColor()).
	//Background(tcell.ColorWhite)

	x, y := GetUpperSlideCorner()
	for idx, row := range slides[currentSlide] {
		DrawStr(screen, x, y+idx, tcell.StyleDefault, row)
	}
	screen.Show()
}

func GetUpperSlideCorner() (int, int) {
	//TODO: range check
	w := len(slides[currentSlide][0])
	for _, row := range slides[currentSlide] {
		if len(row) > w {
			w = len(row)
		}
	}

	h := len(slides[currentSlide])

	sw, sh := screen.Size()

	x := sw/2 - w/2
	y := sh/2 - h/2
	return x, y
}

func Quit() {
	screen.Fini()
	os.Exit(0)
}

func Left() {
	if currentSlide > 0 {
		currentSlide--
	}
}

func Right() {
	if currentSlide < len(slides)-1 {
		currentSlide++
	}
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

	for {
		switch ev := screen.PollEvent().(type) {
		case *tcell.EventResize:
			screen.Sync()
			Draw()
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape:
				Quit()
			case tcell.KeyLeft:
				Left()
			case tcell.KeyRight:
				Right()
			}

			switch ev.Rune() {
			case 'h':
				Left()
			case 'l':
				Right()
			}

			screen.Sync()
			Draw()
		}
	}
}
