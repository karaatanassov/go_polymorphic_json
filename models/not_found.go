package models

import (
	"encoding/json"
	"fmt"

	"gitlab.eng.vmware.com/kkaraatanassov/go-json/interfaces"
)

// NotFound contains the data about a not found error
type NotFound struct {
	RuntimeFault
	ObjKind string
	Obj     string
}

var _ interfaces.NotFound = &NotFound{}
var _ interfaces.RuntimeFault = &NotFound{}
var _ interfaces.Fault = &NotFound{}
var _ json.Marshaler = &NotFound{}
var _ json.Unmarshaler = &NotFound{}

// ZzNotFound is a marker to prevent converting struct with same fields into
// NotFound interface
func (nfo *NotFound) ZzNotFound() {
}

// GetObjKind retrieves the object kind of obj identifier
func (nfo *NotFound) GetObjKind() string {
	return nfo.ObjKind
}

// SetObjKind sets the kind of object references by obj
func (nfo *NotFound) SetObjKind(objKind string) {
	nfo.ObjKind = objKind
}

// GetObj retrieves the obj value
func (nfo *NotFound) GetObj() string {
	return nfo.Obj
}

// SetObj sets the obj id value
func (nfo *NotFound) SetObj(obj string) {
	nfo.Obj = obj
}

// MarshalJSON writes a NotFoundObject as JSON
func (nfo *NotFound) MarshalJSON() ([]byte, error) {
	type marshalable NotFound
	return json.Marshal(struct {
		Kind string
		marshalable
	}{
		Kind:        "NotFound",
		marshalable: marshalable(*nfo),
	})
}

// UnmarshalJSON reads a fault from JSON
func (nfo *NotFound) UnmarshalJSON(in []byte) error {
	pxy := &struct {
		Message string
		Cause   FaultField
		ObjKind string
		Obj     string
	}{}
	err := json.Unmarshal(in, pxy)
	if err != nil {
		return err
	}
	nfo.Message = pxy.Message
	nfo.Cause = pxy.Cause.Fault
	nfo.Obj = pxy.Obj
	nfo.ObjKind = pxy.ObjKind
	return nil
}

// NotFoundField type allows reading polymorphic RuntimeFault fields
type NotFoundField struct {
	interfaces.NotFound
}

var _ interfaces.Fault = &NotFoundField{}
var _ json.Unmarshaler = &NotFoundField{}

// UnmarshalJSON reads the embedded fault taking care of the discriminator
func (ff *NotFoundField) UnmarshalJSON(in []byte) error {
	var err error
	ff.NotFound, err = UnmarshalNotFound(in)
	return err
}

// UnmarshalNotFound reads NotFound or it's subclasses from JSON bytes
func UnmarshalNotFound(in []byte) (interfaces.NotFound, error) {
	fault, err := UnmarshalFault(in)
	if err != nil {
		return nil, err
	}
	if notFound, ok := fault.(interfaces.NotFound); ok {
		return notFound, nil
	}
	return nil, fmt.Errorf("Cannot unmarshal NotFound %v", fault)
}
