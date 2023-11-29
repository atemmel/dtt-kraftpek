package main

import (
	"fmt"
	"strconv"

	"github.com/atemmel/dtt-kraftpek/pkg/md"
	"github.com/gdamore/tcell/v2"
)

var (
	goKeywords = map[string]bool{
		"break":       true,
		"case":        true,
		"chan":        true,
		"const":       true,
		"continue":    true,
		"default":     true,
		"defer":       true,
		"else":        true,
		"fallthrough": true,
		"for":         true,
		"func":        true,
		"go":          true,
		"goto":        true,
		"if":          true,
		"import":      true,
		"interface":   true,
		"map":         true,
		"package":     true,
		"range":       true,
		"return":      true,
		"select":      true,
		"struct":      true,
		"switch":      true,
		"type":        true,
		"var":         true,
	}
)

func isKeyword(word string) bool {
	_, ok := goKeywords[word]
	return ok
}

func isNum(word string) bool {
	_, err := strconv.Atoi(word)
	return err == nil
}

type Renderer struct {
	style      tcell.Style
	Screen     tcell.Screen
	w, h, x, y int
	codeLang   string
}

func NewRenderer(screen tcell.Screen) *Renderer {
	return &Renderer{
		style:  tcell.StyleDefault,
		Screen: screen,
		w:      0,
		h:      0,
		x:      0,
		y:      0,
	}
}

func (r *Renderer) VisitHeader(header *md.Header) {
	r.style = tcell.StyleDefault.
		Foreground(tcell.ColorBlue).
		Background(tcell.ColorBlack).
		Bold(true)

	r.x = r.w - 1

	DrawStr(r.Screen, r.w-2, r.y, r.style, "#")
	header.Child.Accept(r)

	r.style = tcell.StyleDefault
	r.x = r.w
}

func (r *Renderer) VisitCode(code *md.Code) {
	r.codeLang = code.Lang
	for _, child := range code.Children() {
		child.Accept(r)
	}
	r.codeLang = ""
}

func (r *Renderer) VisitList(list *md.List) {
	if list.Ordered {
		r.x = r.w + 2
	} else {
		r.x = r.w + 1
	}

	style := tcell.StyleDefault.
		Foreground(tcell.ColorYellow).
		Background(tcell.ColorBlack).
		Attributes(tcell.AttrBold)

	for i, child := range list.Children() {
		prefix := "●"
		if list.Ordered {
			prefix = strconv.Itoa(i+1) + "."
		}
		DrawStr(r.Screen, r.w, r.y, style, prefix)
		child.Accept(r)
	}

	r.x = r.w
}

func (r *Renderer) VisitText(text *md.Text) {
	switch r.codeLang {
	case "go":
		r.drawColoredCodeText(text.Value)
	default:
		DrawStr(r.Screen, r.x, r.y, r.style, text.Value)
	}
	r.y += 1
}

func (r *Renderer) drawColoredCodeText(text string) {
	wordBegin := 0
	for i := 0; i < len(text); i++ {
		c := text[i]
		if c == '"' {
			strBegin := i
			i++
			for ; i < len(text); i++ {
				if text[i] == '"' {
					break
				}
			}
			str := text[strBegin : i+1]
			prevWord := text[wordBegin:strBegin]
			style := tcell.StyleDefault.
				Foreground(tcell.ColorPurple).
				Background(tcell.ColorBlack)
			DrawStr(r.Screen, r.x+wordBegin, r.y, tcell.StyleDefault, prevWord)
			DrawStr(r.Screen, r.x+strBegin, r.y, style, str)
			wordBegin = i + 1
		} else if i != len(text)-1 && text[i:i+2] == "//" {
			comment := text[i:]
			style := tcell.StyleDefault.
				Foreground(tcell.ColorGray).
				Background(tcell.ColorBlack).
				Attributes(tcell.AttrItalic)
			DrawStr(r.Screen, r.x+wordBegin, r.y, style, comment)
			break
		} else if c == ' ' {
			word := text[wordBegin:i]
			style := tcell.StyleDefault
			if isKeyword(word) {
				style = tcell.StyleDefault.
					Foreground(tcell.ColorYellow).
					Background(tcell.ColorBlack)
			}

			DrawStr(r.Screen, r.x+wordBegin, r.y, style, word)
			wordBegin = i + 1
		} else if i == len(text)-1 {
			word := text[wordBegin : i+1]

			style := tcell.StyleDefault
			if isKeyword(word) {
				style = tcell.StyleDefault.
					Foreground(tcell.ColorYellow).
					Background(tcell.ColorBlack)
			}

			DrawStr(r.Screen, r.x+wordBegin, r.y, style, word)
			wordBegin = i + 1
		}

		if num := text[wordBegin : i+1]; isNum(num) {
			style := tcell.StyleDefault.
				Foreground(tcell.ColorPurple).
				Background(tcell.ColorBlack)
			DrawStr(r.Screen, r.x+wordBegin, r.y, style, num)
			wordBegin = i + 1
		}
	}
}

func (r *Renderer) VisitRoot(root *md.Root) {
	bounds := bounds{}
	r.w, r.h = bounds.Calc(root, r.Screen)
	r.x, r.y = r.w, r.h

	for _, child := range root.Children {
		child.Accept(r)
	}
}

type bounds struct {
	w, h, calcW int
}

func (b *bounds) VisitHeader(header *md.Header) {
	b.calcW = 1
	header.Child.Accept(b)
	b.calcW = 0
}

func (b *bounds) VisitCode(code *md.Code) {
	for _, child := range code.Children() {
		child.Accept(b)
	}
}

func (b *bounds) VisitList(list *md.List) {
	for _, child := range list.Children() {
		child.Accept(b)
	}
}

func (b *bounds) VisitText(text *md.Text) {
	w := len(text.Value) + b.calcW
	if w > b.w {
		b.w = w
	}
	b.h += 1
}

func (b *bounds) VisitRoot(root *md.Root) {
	b.w, b.h = 0, 0
	for _, child := range root.Children {
		child.Accept(b)
	}
}

func (b *bounds) Calc(root *md.Root, screen tcell.Screen) (int, int) {
	b.VisitRoot(root)
	sw, sh := screen.Size()

	x := sw/2 - b.w/2
	y := sh/2 - b.h/2

	return x, y
}

type Printer struct {
	depth int
}

func (p *Printer) up() {
	p.depth -= 2
}

func (p *Printer) down() {
	p.depth += 2
}

func (p *Printer) pad() {
	for i := 0; i < p.depth; i++ {
		fmt.Print(" ")
	}
}

func (p *Printer) VisitHeader(header *md.Header) {
	p.down()
	p.pad()
	fmt.Println("Header")
	header.Child.Accept(p)
	p.up()
}

func (p *Printer) VisitCode(code *md.Code) {
	p.down()
	p.pad()
	fmt.Printf("Code, lang: %s\n", code.Lang)
	for _, child := range code.Children() {
		child.Accept(p)
	}
	p.up()
}

func (p *Printer) VisitList(list *md.List) {
	p.down()
	p.pad()
	fmt.Println("List")
	for _, child := range list.Children() {
		child.Accept(p)
	}
	p.up()
}

func (p *Printer) VisitText(text *md.Text) {
	p.down()
	p.pad()
	fmt.Println("Text:", text.Value)
	p.up()
}

func (p *Printer) VisitRoot(root *md.Root) {
	fmt.Println("Root")
	for _, child := range root.Children {
		child.Accept(p)
	}
}
