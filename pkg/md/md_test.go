package md

import "testing"

func TestMyFunc(t *testing.T) {
	if !MyFunc() {
		t.Fatal("Det small i MyFunc")
	}
}
