package main

import (
	"errors"
	"fmt"

	"github.com/atemmel/dtt-kraftpek/pkg/md"
	"github.com/gdamore/tcell/v2"
)

type Renderer struct {
	style      tcell.Style
	Screen     tcell.Screen
	w, h, x, y int
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
		Foreground(tcell.ColorBlue.TrueColor()).
		Bold(true)

	r.x = r.w - 1

	DrawStr(r.Screen, r.w-2, r.y, r.style, "#")
	header.Child.Accept(r)

	r.style = tcell.StyleDefault
	r.x = r.w
}

func (r *Renderer) VisitCode(code *md.Code) {
	panic(errors.New("Ohanterad nod"))
}

func (r *Renderer) VisitList(list *md.List) {
	r.x = r.w - 1

	for _, child := range list.Children() {
		DrawStr(r.Screen, r.w-2, r.y, r.style, "●")
		child.Accept(r)
	}

	r.x = r.w
}

func (r *Renderer) VisitText(text *md.Text) {
	DrawStr(r.Screen, r.x, r.y, r.style, text.Value)
	r.y += 1
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
	panic(errors.New("Ohanterad nod"))
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
	panic(errors.New("Ohanterad nod"))
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
