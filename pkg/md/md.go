package md

import (
	"strings"
	"unicode"
)

type Node interface {
	Children() []Node
}

type Header struct {
	Child Node
	Level int
}

type Code struct {
	Child []Node
	Lang  string
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
	return l.Children()
}
func (c *Code) Children() []Node {
	return c.Children()
}

func ParseMd(src string) Root {
	src = strings.TrimSpace(src)
	root := Root{}

	for i := 0; i < len(src); i++ {

		//if child := readCode(&i, src); child != nil {
		//	root.Children = append(root.Children, child)
		//	continue
		//}

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

	if len(src)-*index < 3 {
		return nil
	}

	if src[*index:3] != "```" {
		return nil
	}
	*index += 3
	language := ""
	textBegin := *index
	for ; *index < len(src); *index++ {

		if src[*index] == '\n' {
			language = src[textBegin:*index]
		}

	}

	print(readText(index, src))

	return &Code{
		Lang: language,
	}

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
		if src[*index+1] != '.' {
			*index--
			break
		}
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
		list.children = append(list.children, readText(index, src))
	}

	return list
}

func readText(index *int, src string) Node {
	textBegin := *index
	for ; *index < len(src); *index++ {
		c := src[*index]
		if c == '\n' {
			line := src[textBegin:*index]
			textBegin = *index + 1

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
