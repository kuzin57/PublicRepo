package writer

type IWriter interface {
	Write(v any) error
}
