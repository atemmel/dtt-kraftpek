package md

import (
	"strings"
	"unicode"
)

type Visitor interface {
	VisitHeader(*Header)
	VisitCode(*Code)
	VisitList(*List)
	VisitText(*Text)
	VisitRoot(*Root)
}

type Node interface {
	Children() []Node
	Accept(Visitor)
}

type Header struct {
	Child Node
	Level int
}

type Code struct {
	children []Node
	Lang     string
}

type List struct {
	children []Node
	Ordered  bool
}

type Text struct {
	Value string
}

type Root struct {
	Children []Node
}

func (_ Text) Children() []Node {
	return []Node{}
}

func (h *Header) Children() []Node {
	return []Node{h.Child}
}

func (l *List) Children() []Node {
	return l.children
}

func (c *Code) Children() []Node {
	return c.children
}

func (t *Text) Accept(v Visitor) {
	v.VisitText(t)
}

func (h *Header) Accept(v Visitor) {
	v.VisitHeader(h)
}

func (l *List) Accept(v Visitor) {
	v.VisitList(l)
}

func (c *Code) Accept(v Visitor) {
	v.VisitCode(c)
}

func ParseMd(src string) Root {
	src = strings.TrimSpace(src)
	root := Root{}

	for i := 0; i < len(src); i++ {

		if child := readCode(&i, src); child != nil {
			root.Children = append(root.Children, child)
			continue
		}

		if child := readUnorderedList(&i, src); child != nil {
			root.Children = append(root.Children, child)
			continue
		}
		if child := readOrderedList(&i, src); child != nil {
			root.Children = append(root.Children, child)
			continue
		}

		if child := readHeader(&i, src); child != nil {
			root.Children = append(root.Children, child)
			continue
		}

		if child := readText(&i, src); child != nil {
			root.Children = append(root.Children, child)
			continue
		}

	}

	return root
}

func readCode(index *int, src string) Node {
	tempindex := *index
	if len(src)-*index < 3 {
		return nil
	}

	if src[*index:*index+3] != "```" {
		return nil
	}
	*index += 3

	code := &Code{
		children: make([]Node, 0),
		Lang:     "",
	}

	for langBegin := *index; *index < len(src); *index++ {
		if src[*index] == '\n' {
			code.Lang = src[langBegin:*index]
			*index++
			break
		}
	}

	for ; *index < len(src); *index++ {
		if *index+1 >= len(src) {
			*index = tempindex
			return nil
		}

		c := src[*index : *index+3]
		if c == "```" {
			*index += 2
			break
		}
		code.children = append(code.children, readText(index, src))
	}

	return code

}

func readHeader(index *int, src string) Node {

	c := src[*index]

	if c != '#' {
		return nil
	}
	*index++

	return &Header{
		Child: readText(index, src),
		Level: 1,
	}

}

func readOrderedList(index *int, src string) Node {
	if !unicode.IsNumber(rune(src[*index])) {
		return nil
	}
	if src[*index+1] != '.' {
		return nil
	}

	list := &List{
		children: make([]Node, 0),
		Ordered:  true,
	}

	for ; *index < len(src); *index++ {
		if !unicode.IsNumber(rune(src[*index])) {
			*index--
			break
		}
		if *index+1 >= len(src) {
			*index--
			break
		}

		if src[*index+1] != '.' {
			*index--
			break
		}
		*index += 2
		list.children = append(list.children, readText(index, src))
	}

	return list
}

func readUnorderedList(index *int, src string) Node {
	if src[*index] != '*' {
		return nil
	}

	list := &List{
		children: make([]Node, 0),
		Ordered:  false,
	}

	for ; *index < len(src); *index++ {
		c := src[*index]
		if c != '*' {
			*index--
			break
		}
		*index++
		list.children = append(list.children, readText(index, src))
	}

	return list
}

func readText(index *int, src string) Node {
	textBegin := *index
	for ; *index < len(src); *index++ {
		if src[*index] == '\n' {
			line := src[textBegin:*index]

			return &Text{
				Value: line,
			}
		}

	}
	line := src[textBegin:]

	return &Text{
		Value: line,
	}
}
