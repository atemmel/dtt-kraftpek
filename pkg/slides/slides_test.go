package slides

import (
	"testing"
)

func TestReadKraftfil(t *testing.T) {
	files, err := readKraftfil("../../samples/kraftfil")
	if err != nil {
		t.Fatal("Fel uppstod vid läsning av kraftfil:", err)
	}
	if len(files) != 1 {
		t.Fatal("Förväntades läsa 1 fil, läste:", len(files))
	}

	if files[0] != "intro.md" {
		t.Fatal("Förväntades hitta 'intro.md', hittade:", files[0])
	}
}

func TestReadSlides(t *testing.T) {
	slides, err := ReadSlides("../../samples")
	if err != nil {
		t.Fatal("Fel uppstod vid läsning av slides:", err)

	}

	if len(slides) != 1 {
		t.Fatal("Förväntades läsa 1 slide, läste:", len(slides))
	}
}
