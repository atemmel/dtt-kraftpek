package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/atemmel/dtt-kraftpek/pkg/md"
	"github.com/atemmel/dtt-kraftpek/pkg/slides"
	"github.com/gdamore/tcell/v2"
	_ "github.com/gdamore/tcell/v2/encoding"
)

var (
	debugFlag = flag.String("debug", "", "Felsök syntaxträd av en markdownfil")
	slidesDir = "."

	currentSlide = 0
	loadedSlides []slides.Slide
)

func Draw(screen tcell.Screen, renderer *Renderer) {
	screen.Clear()
	// exempel på hur styling kan sättas

	//style := tcell.StyleDefault.
	//Foreground(tcell.ColorCadetBlue.TrueColor()).
	//Background(tcell.ColorWhite)

	root := &loadedSlides[currentSlide].Root
	renderer.VisitRoot(root)

	screen.Show()
}

func Quit(screen tcell.Screen) {
	screen.Fini()
	os.Exit(0)
}

func Left() {
	if currentSlide > 0 {
		currentSlide--
	}
}

func Right() {
	if currentSlide < len(loadedSlides)-1 {
		currentSlide++
	}
}

func debug() {
	bytes, err := os.ReadFile(*debugFlag)
	if err != nil {
		fmt.Println("Kunde inte felsöka fil:", debugFlag)
		os.Exit(2)
	}

	root := md.ParseMd(string(bytes))
	printer := Printer{}
	printer.VisitRoot(&root)
}

func init() {
	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		slidesDir = args[0]
	}
}

func main() {
	var err error

	if *debugFlag != "" {
		debug()
		return
	}

	loadedSlides, err = slides.ReadSlides(slidesDir)
	if err != nil {
		fmt.Println("Kunde inte hitta slides i:", slidesDir)
		os.Exit(1)
	}

	screen, err := tcell.NewScreen()
	Assert(err)
	defer TuiPanic(screen)
	Assert(screen.Init())

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)

	screen.SetStyle(defStyle)

	renderer := NewRenderer(screen)

	for {
		Draw(screen, renderer)
		switch ev := screen.PollEvent().(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape:
				Quit(screen)
			case tcell.KeyLeft:
				Left()
			case tcell.KeyRight:
				Right()
			}

			switch ev.Rune() {
			case 'q':
				Quit(screen)
			case 'h':
				Left()
			case 'l':
				Right()
			}
		}
	}
}
