package csv

import (
	"io"

	"gitlab.com/slon/shad-go/gitfame/pkg/gocsv"
	"gitlab.com/slon/shad-go/gitfame/pkg/writer"
)

type csvWriter struct {
	base io.Writer
}

func NewWriter(base io.Writer) writer.IWriter {
	return &csvWriter{
		base: base,
	}
}

func (w *csvWriter) Write(obj any) error {
	return gocsv.Marshal(obj, w.base)
}
