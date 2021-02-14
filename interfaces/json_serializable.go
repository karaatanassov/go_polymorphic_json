package interfaces

// JSONSerializable is utility interface that avoids infinite recursion when
// writing discriminator using small utlity class. It is possible to overcome
// infinite recursion by coying all of the object fields
type JSONSerializable interface {
	SerializeJSON() ([]byte, error)
}
