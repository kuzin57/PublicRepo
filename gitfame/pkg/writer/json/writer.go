package json

import (
	"encoding/json"
	"io"

	"gitlab.com/slon/shad-go/gitfame/pkg/writer"
)

type jsonWriter struct {
	base io.Writer
}

func NewWriter(base io.Writer) writer.IWriter {
	return &jsonWriter{
		base: base,
	}
}

func (w *jsonWriter) Write(obj any) error {
	content, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	_, err = w.base.Write(content)
	return err
}
