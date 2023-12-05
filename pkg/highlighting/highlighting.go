package highlighting

import "strings"

type Kind int

const (
	Normal Kind = iota
	Type
	Keyword
	StringLiteral
	NumberLiteral
	Value
	Function
	Comment
)

type Fragment struct {
	Kind
	Content string
}

func Parse(src string) []Fragment {
	src = strings.TrimRight(src, "\n")
	fragments := make([]Fragment, 0, 4)

	fragmentBegin := 0

	for i := 0; i < len(src); i++ {
		c := src[i]
		content := src[fragmentBegin:i]

		if content == " " {
			fragments = append(fragments, Fragment{
				Content: content,
				Kind:    Normal,
			})
			fragmentBegin = i
		} else if content == "\t" {
			fragments = append(fragments, Fragment{
				Content: "    ",
				Kind:    Normal,
			})
			fragmentBegin = i
		} else if content == "//" {
			fragments = append(fragments, Fragment{
				Content: src[fragmentBegin:],
				Kind:    Comment,
			})
			fragmentBegin = len(src)
			break
		} else if content == "\"" {
			for ; i < len(src); i++ {
				if src[i] == '"' {
					i++
					break
				}
			}

			fragments = append(fragments, Fragment{
				Content: src[fragmentBegin:i],
				Kind:    StringLiteral,
			})
			fragmentBegin = i
		} else if isOnlySymbols(content) && !isSymbol(c) {
			if content == "" {
				continue
			}

			fragments = append(fragments, Fragment{
				Content: content,
				Kind:    Normal,
			})

			fragmentBegin = i
		} else if c == ' ' || (!isOnlySymbols(content) && isSymbol(c)) {
			if content == "" {
				continue
			}

			kind := lookupFragmentKindFromContent(content)

			fragments = append(fragments, Fragment{
				Content: content,
				Kind:    kind,
			})

			fragmentBegin = i
		}
	}

	content := src[fragmentBegin:]
	if content != "" {
		kind := lookupFragmentKindFromContent(content)

		fragments = append(fragments, Fragment{
			Content: content,
			Kind:    kind,
		})
	}

	return fragments
}

func lookupFragmentKindFromContent(content string) Kind {
	if isKeyword(content) {
		return Keyword
	}

	if isType(content) {
		return Type
	}

	if isNum(content) {
		return NumberLiteral
	}

	if isValue(content) {
		return Value
	}

	if isFn(content) {
		return Function
	}

	return Normal
}
