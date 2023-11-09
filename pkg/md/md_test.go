package md

import (
	"testing"
)

func TestParseMd(t *testing.T) {
	src := `hej
på
dig
`

	root := ParseMd(src)

	if len(root.Children) != 3 {
		t.Fatal("Förväntades hitta 3 barn, hittade", len(root.Children))
	}

	if root.Children[0].(*Text).Value != "hej" {
		t.Fatal("Förväntades hitta hej, hittade istället : ", root.Children[0].(*Text).Value)

	}

}

func TestParseMdHeader(t *testing.T) {
	src := `#hej
på
dig
`

	root := ParseMd(src)

	if len(root.Children) != 3 {
		t.Fatal("Förväntades hitta 3 barn, hittade", len(root.Children))
	}

	if root.Children[0].(*Header).Level != 1 {
		t.Fatal("Förväntades level är 1, hittade istället : ", root.Children[0].(*Header).Level)

	}
}

func TestParseMdUnorderedList(t *testing.T) {
	src := `hej
*på
*dig
Jonas
`

	root := ParseMd(src)

	if len(root.Children) != 3 {
		t.Fatal("Förväntades hitta 3 barn, hittade", len(root.Children))
	}

	if root.Children[1].(*List).Ordered {
		t.Fatal("Förväntades hitta oordnad lista")
	}
}

func TestParseMdOrderedList(t *testing.T) {
	src := `hej
1.på
2.dig
Jonas
`

	root := ParseMd(src)

	if len(root.Children) != 3 {
		t.Fatal("Förväntades hitta 3 barn, hittade", len(root.Children))
	}

	if !root.Children[1].(*List).Ordered {
		t.Fatal("Förväntades hitta ordnad lista")
	}
}

func TestParseMdbadalistor(t *testing.T) {
	src := `#en header
1.object 1
2.okbject 2
* uobject 1
* uobject 2
random text
`

	root := ParseMd(src)

	if len(root.Children) != 4 {
		t.Fatal("Förväntades hitta 4 barn, hittade", len(root.Children))
	}
	if root.Children[0].(*Header).Level != 1 {
		t.Fatal("Förväntades hitta header level1")
	}
	if !root.Children[1].(*List).Ordered {
		t.Fatal("Förväntades hitta ordnad lista")
	}
	if root.Children[1].(*List).Ordered {
		t.Fatal("Förväntades hitta oordnad lista")
	}
}
