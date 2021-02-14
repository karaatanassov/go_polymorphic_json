package interfaces

// Fault represents a base error
// To be generated
type Fault interface {
	JSONSerializable
	GetMessage() string
	SetMessage(string)
}
