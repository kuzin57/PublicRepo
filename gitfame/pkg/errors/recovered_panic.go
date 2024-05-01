package errors

import (
	"fmt"
)

type RecoveredPanicError struct {
	panicPayload interface{}
	stacktrace   []byte
}

func NewRecoveredPanicError(panicPayload interface{}, stacktrace []byte) *RecoveredPanicError {
	fmt.Println("Panic!", panicPayload, string(stacktrace))
	if panicPayload == nil {
		return nil
	}
	return &RecoveredPanicError{
		panicPayload: panicPayload,
		stacktrace:   stacktrace,
	}
}

func (e *RecoveredPanicError) Error() string {
	return fmt.Sprintf("Panic handeled\n%v\n%s", e.panicPayload, e.stacktrace)
}

func (e *RecoveredPanicError) Unwrap() error {
	if err, ok := e.panicPayload.(error); ok {
		return err
	}
	return nil
}
