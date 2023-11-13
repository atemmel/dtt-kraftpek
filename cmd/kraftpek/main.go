package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/atemmel/dtt-kraftpek/pkg/md"
	"github.com/atemmel/dtt-kraftpek/pkg/slides"
	"github.com/gdamore/tcell/v2"
	_ "github.com/gdamore/tcell/v2/encoding"
	"github.com/mattn/go-runewidth"
)

var (
	screen tcell.Screen

	currentSlide = 0
	Slides []slides.Slide
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

func visitRoot(root *md.Root, callback func(md.Node)) {
	for _, child := range root.Children {
		visit(child, callback)
	}
}

func visit(node md.Node, callback func(md.Node)) {
	callback(node)
	for _, child := range node.Children() {
		visit(child, callback)
	}
}

func DrawRoot(screen tcell.Screen, root *md.Root) {
	w, h := GetUpperSlideCorner(root)
	x, y := w, h

	writeCallback := func(n md.Node) {
		switch node := n.(type) {
		case *md.Text:
			DrawStr(screen, x, y, tcell.StyleDefault, node.Value)
			y++
		default:
			panic(errors.New("Okänd nod funnen i markdown " + reflect.TypeOf(node).String()))
		}
	}

	visitRoot(root, writeCallback)
}

func getRenderSize(root *md.Root) (int, int) {
	maxWidth := 0
	maxHeight := 0

	widthCallback := func(n md.Node) {
		switch node := n.(type) {
		case *md.Text:
			if len(node.Value) > maxWidth {
				maxWidth = len(node.Value)
			}
			maxHeight += 1
		default:
			panic(errors.New("Okänd nod funnen i markdown " + reflect.TypeOf(node).String()))
		}
	}

	visitRoot(root, widthCallback)

	return maxWidth, maxHeight
}

func Draw() {
	screen.Clear()
	// exempel på hur styling kan sättas

	//style := tcell.StyleDefault.
	//Foreground(tcell.ColorCadetBlue.TrueColor()).
	//Background(tcell.ColorWhite)

	root := &Slides[currentSlide].Root
	DrawRoot(screen, root)

	screen.Show()
}

func GetUpperSlideCorner(root *md.Root) (int, int) {
	w, h := getRenderSize(root)
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
	if currentSlide < len(Slides)-1 {
		currentSlide++
	}
}

var (
	slidesDir = "."
)

func init() {
	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		slidesDir = args[0]
	}
}

func main() {
	var err error
	Slides, err = slides.ReadSlides(slidesDir)
	if err != nil {
		fmt.Println("Kunde inte hitta slides i:", slidesDir)
		os.Exit(1)
	}

	screen, err = tcell.NewScreen()
	Assert(err)
    defer func() {
        if r := recover(); r != nil {
			screen.Fini()
			fmt.Println(r)
        }
    }()
	err = screen.Init()
	Assert(err)

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)

	for {
		Draw()
		switch ev := screen.PollEvent().(type) {
		case *tcell.EventResize:
			screen.Sync()
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
			case 'q':
				Quit()
			case 'h':
				Left()
			case 'l':
				Right()
			}
		}
	}
}
