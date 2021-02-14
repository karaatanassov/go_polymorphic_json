package interfaces

// NotFound represents error when object is not found
// This should be generated code. It is possible ot emit it in the models
// package and creative name will be needed e.g. NotFoundInterface.
// The interface package approach seems cleaner.
type NotFound interface {
	RuntimeFault
	GetObjKind() string
	SetObjKind(string)
	GetObj() string
	SetObj(string)
}
