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

func TestParseMdList(t *testing.T) {
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
