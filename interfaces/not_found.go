package interfaces

// NotFound represents error when object is not found
// To be generated
type NotFound interface {
	RuntimeFault
	GetObjKind() string
	SetObjKind(string)
	GetObj() string
	SetObj(string)
	ZzNotFound()
}
