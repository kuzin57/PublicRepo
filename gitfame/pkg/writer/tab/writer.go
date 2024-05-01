package tab

import (
	"io"
	"reflect"
	"strconv"

	"gitlab.com/slon/shad-go/gitfame/pkg/writer"
)

const (
	EOL = '\n'
)

type tabWriter struct {
	base io.Writer
	tabs []int
}

func NewWriter(base io.Writer) writer.IWriter {
	return &tabWriter{
		base: base,
	}
}

func (w *tabWriter) writeOne(v reflect.Value) error {
	numFields := v.NumField()

	var index int
	for i := range numFields {
		field := v.Field(i)

		toWrite := w.getFieldRepresentation(field)
		if toWrite == "" {
			continue
		}

		if _, err := w.base.Write([]byte(toWrite)); err != nil {
			return err
		}

		if i < numFields-1 {
			if err := w.writeSpaces(len(toWrite), w.tabs[index]); err != nil {
				return err
			}
		} else {
			if _, err := w.base.Write([]byte{EOL}); err != nil {
				return err
			}
		}
		index++
	}

	return nil
}

func (w *tabWriter) writeSpaces(fromLen, toLen int) error {
	for i := fromLen; i < toLen; i++ {
		if _, err := w.base.Write([]byte{' '}); err != nil {
			return err
		}
	}

	return nil
}

func (w *tabWriter) analyzeTabs(v reflect.Value) {
	firstElem := v.Index(0)
	for i := range firstElem.NumField() {
		tag := firstElem.Type().Field(i).Tag.Get("tab")
		if tag == "-" {
			continue
		}

		maxLen := len(tag) + 1
		w.tabs = append(w.tabs, maxLen)
	}

	for i := range v.Len() {
		val := v.Index(i)

		var index int
		for j := range val.NumField() {
			strRepr := w.getFieldRepresentation(val.Field(j))
			if strRepr == "" {
				continue
			}

			if len(strRepr) > w.tabs[index]-1 {
				w.tabs[index] = len(strRepr) + 1
			}
			index++
		}
	}
}

func (w *tabWriter) getFieldRepresentation(field reflect.Value) string {
	switch {
	case field.Kind() == reflect.String:
		return field.String()

	case field.Kind() == reflect.Int:
		return strconv.Itoa(int(field.Int()))

	default:
		return ""
	}
}

func (w *tabWriter) writeHeadLine(v reflect.Value) error {
	var index int
	for i := range v.NumField() {
		tag := v.Type().Field(i).Tag.Get("tab")
		if tag == "-" {
			continue
		}

		if _, err := w.base.Write([]byte(tag)); err != nil {
			return err
		}

		if i < v.NumField()-1 {
			if err := w.writeSpaces(len(tag), w.tabs[index]); err != nil {
				return err
			}
		}

		index++
	}

	_, err := w.base.Write([]byte{EOL})
	return err
}

// write slices beautifully
func (w *tabWriter) Write(obj any) error {
	v := reflect.ValueOf(obj)
	w.analyzeTabs(v)

	for i := range v.Len() {
		if i == 0 {
			if err := w.writeHeadLine(v.Index(i)); err != nil {
				return err
			}
		}

		if err := w.writeOne(v.Index(i)); err != nil {
			return err
		}
	}

	return nil
}
