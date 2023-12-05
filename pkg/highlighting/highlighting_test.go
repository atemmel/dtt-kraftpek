package highlighting

import "testing"

func AssertFragmentLength(t *testing.T, fragments []Fragment, expected int) {
	if len(fragments) != expected {
		t.Fatalf("Unexpected fragment length, expected '%d', got '%d'", expected, len(fragments))
	}
}

func AssertFragmentKind(t *testing.T, fragment Fragment, expected Kind) {
	if fragment.Kind != expected {
		t.Errorf("Unexpected fragment kind, expected '%d', got '%d'", expected, fragment.Kind)
	}
}

func AssertFragmentContent(t *testing.T, fragment Fragment, expected string) {
	if fragment.Content != expected {
		t.Errorf("Unexpected fragment content, expected '%s', got '%s'", expected, fragment.Content)
	}
}

func TestParseNormal(t *testing.T) {
	src := "hi hello hi"

	fragments := Parse(src)

	AssertFragmentLength(t, fragments, 5)

	AssertFragmentKind(t, fragments[0], Normal)

	AssertFragmentContent(t, fragments[0], "hi")
}

func TestParseKeyword(t *testing.T) {
	src := "for x := range y {"

	fragments := Parse(src)

	AssertFragmentContent(t, fragments[0], "for")
	AssertFragmentKind(t, fragments[0], Keyword)
	AssertFragmentContent(t, fragments[1], " ")
	AssertFragmentContent(t, fragments[2], "x")
	AssertFragmentContent(t, fragments[3], " ")
	AssertFragmentContent(t, fragments[4], ":=")
	AssertFragmentContent(t, fragments[5], " ")
	AssertFragmentContent(t, fragments[6], "range")
	AssertFragmentKind(t, fragments[6], Keyword)
}

func TestParseType(t *testing.T) {
	src := "func(x string) (int, int)"

	fragments := Parse(src)

	AssertFragmentContent(t, fragments[0], "func")
	AssertFragmentKind(t, fragments[0], Keyword)
	AssertFragmentContent(t, fragments[1], "(")
	AssertFragmentContent(t, fragments[2], "x")
	AssertFragmentContent(t, fragments[3], " ")
	AssertFragmentContent(t, fragments[4], "string")
	AssertFragmentKind(t, fragments[4], Type)
	AssertFragmentContent(t, fragments[5], ")")
	AssertFragmentContent(t, fragments[6], " ")
	AssertFragmentContent(t, fragments[7], "(")
	AssertFragmentContent(t, fragments[8], "int")
	AssertFragmentKind(t, fragments[8], Type)
	AssertFragmentContent(t, fragments[9], ",")
	AssertFragmentContent(t, fragments[10], " ")
	AssertFragmentContent(t, fragments[11], "int")
	AssertFragmentKind(t, fragments[11], Type)
	AssertFragmentContent(t, fragments[12], ")")
}

func TestParseStringLiteral(t *testing.T) {
	src := "\"text\" for \"days\""

	fragments := Parse(src)

	AssertFragmentContent(t, fragments[0], "\"text\"")
	AssertFragmentKind(t, fragments[0], StringLiteral)
	AssertFragmentContent(t, fragments[1], " ")
	AssertFragmentContent(t, fragments[2], "for")
	AssertFragmentContent(t, fragments[3], " ")
	AssertFragmentContent(t, fragments[4], "\"days\"")
	AssertFragmentKind(t, fragments[4], StringLiteral)
}

func TestParseNumberLiteral(t *testing.T) {
	src := "1 by 1"

	fragments := Parse(src)

	AssertFragmentContent(t, fragments[0], "1")
	AssertFragmentKind(t, fragments[0], NumberLiteral)
	AssertFragmentContent(t, fragments[1], " ")
	AssertFragmentContent(t, fragments[2], "by")
	AssertFragmentContent(t, fragments[3], " ")
	AssertFragmentContent(t, fragments[4], "1")
	AssertFragmentKind(t, fragments[4], NumberLiteral)
}

func TestParseComment(t *testing.T) {
	src := "lol // this is comment"

	fragments := Parse(src)

	AssertFragmentContent(t, fragments[0], "lol")
	AssertFragmentContent(t, fragments[1], " ")
	AssertFragmentContent(t, fragments[2], "// this is comment")
	AssertFragmentKind(t, fragments[2], Comment)
}

func TestParseFailure1(t *testing.T) {
	src := "make(map[string]int)"

	fragments := Parse(src)

	AssertFragmentLength(t, fragments, 8)
	AssertFragmentContent(t, fragments[0], "make")
	AssertFragmentContent(t, fragments[1], "(")
	AssertFragmentContent(t, fragments[2], "map")
}

func TestParseFailure2(t *testing.T) {
	src := "fmt.Println(\"j är\", j)   // skriv ut det nya värdet på j"

	fragments := Parse(src)

	AssertFragmentContent(t, fragments[0], "fmt.Println")
	AssertFragmentContent(t, fragments[1], "(")
	AssertFragmentContent(t, fragments[2], "\"j är\"")
	AssertFragmentContent(t, fragments[3], ",")
	AssertFragmentContent(t, fragments[4], " ")
	AssertFragmentContent(t, fragments[5], "j")
	AssertFragmentContent(t, fragments[6], ")")
	AssertFragmentContent(t, fragments[7], " ")
	AssertFragmentContent(t, fragments[8], " ")
	AssertFragmentContent(t, fragments[9], " ")
	AssertFragmentContent(t, fragments[10], "// skriv ut det nya värdet på j")
}

func TestParseFailure3(t *testing.T) {
	src := "if err != nil {"

	fragments := Parse(src)

	AssertFragmentContent(t, fragments[0], "if")
	AssertFragmentContent(t, fragments[1], " ")
	AssertFragmentContent(t, fragments[2], "err")
	AssertFragmentContent(t, fragments[3], " ")
	AssertFragmentContent(t, fragments[4], "!")
	AssertFragmentContent(t, fragments[5], "=")
	AssertFragmentContent(t, fragments[6], " ")
	AssertFragmentContent(t, fragments[7], "nil")

}
