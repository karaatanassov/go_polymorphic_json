package interfaces

// Fault represents a base error
// This should be generated code. It is possible to emit it in the models
// package and creative name will be needed e.g. NotFoundInterface.
// The interface package approach seems cleaner.
type Fault interface {
	GetMessage() string
	SetMessage(string)
}
