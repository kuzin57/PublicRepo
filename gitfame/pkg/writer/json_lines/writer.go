package jsonlines

import (
	"encoding/json"
	"io"
	"reflect"

	"gitlab.com/slon/shad-go/gitfame/pkg/writer"
)

var EOL = []byte{'\n'}

type jsonLWriter struct {
	base io.Writer
}

func NewWriter(base io.Writer) writer.IWriter {
	return &jsonLWriter{
		base: base,
	}
}

func (w *jsonLWriter) writeOne(v any) error {
	content, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = w.base.Write(content)
	return err
}

func (w *jsonLWriter) Write(obj any) error {
	if reflect.TypeOf(obj).Kind() != reflect.Slice {
		return w.writeOne(obj)
	}

	v := reflect.ValueOf(obj)
	for i := range v.Len() {
		if err := w.writeOne(v.Index(i).Interface()); err != nil {
			return err
		}

		if _, err := w.base.Write(EOL); err != nil {
			return err
		}
	}

	return nil
}
