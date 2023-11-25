package merger

import (
	"bytes"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"io"
)

func MergeTwoPages(pageOne, pageTwo io.ReadSeeker) (merged *bytes.Buffer, err error) {
	err = api.MergeRaw([]io.ReadSeeker{pageOne, pageTwo}, merged, nil)
	return
}
