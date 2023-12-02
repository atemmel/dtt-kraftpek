package highlighting

import "strconv"

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

	goTypes = map[string]bool{
		"string": true,
		"int":    true,
		"float":  true,
		"bool":   true,
	}
)

func isKeyword(word string) bool {
	_, ok := goKeywords[word]
	return ok
}

func isType(word string) bool {
	_, ok := goTypes[word]
	return ok
}

func isSymbol(prospect byte) bool {
	switch prospect {
	case '(', ')', ',', '{', '}', '=', '<', '>', ':', '[', ']', '&', '*', ';':
		return true
	}
	return false
}

func isOnlySymbols(prospect string) bool {
	for i := range prospect {
		c := prospect[i]
		if !isSymbol(c) {
			return false
		}
	}
	return true
}

func isNum(word string) bool {
	_, err := strconv.Atoi(word)
	return err == nil
}
