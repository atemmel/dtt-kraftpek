package slides

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/atemmel/dtt-kraftpek/pkg/md"
)

type Slide struct {
	Name string
	Root md.Root
}

type readResult struct {
	Slide
	Number int
	Error error
}

func readKraftfil(where string) ([]string, error) {
	bytes, err := os.ReadFile(where)
	if err != nil {
		return nil, err
	}

	clean := strings.TrimSpace(string(bytes))
	return strings.Split(clean, "\n"), nil
}

func ReadSlides(where string) ([]Slide, error) {
	files, err := readKraftfil(path.Join(where, "kraftfil"))
	if err != nil {
		return nil, err
	}

	channel := make(chan readResult)
	for i, file := range files {
		fullpath := path.Join(where, file)
		go readSlide(channel, fullpath, i)
	}

	slides := make([]Slide, len(files))

	for i := 0; i < len(files); i++ {
		result := <- channel
		if result.Error != nil {
			return nil, result.Error
		}

		slides[result.Number] = result.Slide
	}


	return slides, nil
}

func readSlide(channel chan readResult, which string, number int) {
	name := path.Base(which)
	bytes, err := os.ReadFile(which)
	if err != nil {
		channel <- readResult{
			Slide: Slide{},
			Number: number,
			Error: errors.New("Kunde inte läsa fil: " + which + err.Error()),
		}
	}

	root := md.ParseMd(string(bytes))

	channel <- readResult{
		Slide: Slide{
			Name: name,
			Root: root,
		},
		Number: number,
		Error: nil,
	}
}
