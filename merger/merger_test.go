package merger

import (
	"os"
	"testing"
)

func TestMerger(t *testing.T) {
	f1, err := os.Open("output_4194006.pdf")
	if err != nil {
		t.Error(err)
	}
	defer f1.Close()

	f2, err := os.Open("output2_4194006.pdf")
	if err != nil {
		t.Error(err)
	}
	defer f2.Close()

	w, err := MergeTwoPages(f1, f2)
	if err != nil {
		t.Error(err)
	}

	f, err := os.Create("merged.pdf")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	_, err = f.Write(w.Bytes())
	if err != nil {
		t.Error(err)
	}
}
