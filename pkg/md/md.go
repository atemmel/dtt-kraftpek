package md

import "strings"

type Node interface {
	Children() []Node
}

type Header struct {
	Child Node
	Level int
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

func ParseMd(src string) Root {
	src = strings.TrimSpace(src)
	root := Root{}

	for i := 0; i < len(src); i++ {

		if child := readHeader(&i, src); child != nil {
			root.Children = append(root.Children, child)
			continue
		}

		if child := readList(&i, src); child != nil {
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

func readList(index *int, src string) Node {

	return readUnorderedList(index, src)
	//readOrderedList()

}

func readOrderedList(index *int, src string) Node {
	return nil
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
