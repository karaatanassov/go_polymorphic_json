package raw_message

import (
	"encoding/json"
	"fmt"
)

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

// NotFoundStruct contains the data about a not found error
type NotFoundStruct struct {
	RuntimeFaultStruct
	ObjKind string
	Obj     string
}

var _ NotFound = &NotFoundStruct{}
var _ RuntimeFault = &NotFoundStruct{}
var _ Fault = &NotFoundStruct{}
var _ json.Marshaler = &NotFoundStruct{}
var _ json.Unmarshaler = &NotFoundStruct{}

// ZzNotFound is a marker to prevent converting struct with same fields into
// NotFound interface
func (nfo *NotFoundStruct) ZzNotFound() {
}

// GetObjKind retrieves the object kind of obj identifier
func (nfo *NotFoundStruct) GetObjKind() string {
	return nfo.ObjKind
}

// SetObjKind sets the kind of object references by obj
func (nfo *NotFoundStruct) SetObjKind(objKind string) {
	nfo.ObjKind = objKind
}

// GetObj retrieves the obj value
func (nfo *NotFoundStruct) GetObj() string {
	return nfo.Obj
}

// SetObj sets the obj id value
func (nfo *NotFoundStruct) SetObj(obj string) {
	nfo.Obj = obj
}

// MarshalJSON writes a NotFoundObject as JSON
func (nfo *NotFoundStruct) MarshalJSON() ([]byte, error) {
	type marshalable NotFoundStruct
	return json.Marshal(struct {
		Kind string
		marshalable
	}{
		Kind:        "NotFound",
		marshalable: marshalable(*nfo),
	})
}

// UnmarshalJSON reads a fault from JSON
func (nfo *NotFoundStruct) UnmarshalJSON(in []byte) error {
	pxy := &struct {
		Message string
		Cause   json.RawMessage
		ObjKind string
		Obj     string
	}{}
	err := json.Unmarshal(in, pxy)
	if err != nil {
		return err
	}
	var cause Fault
	if pxy.Cause != nil {
		cause, err = UnmarshalFault(pxy.Cause)
		if err != nil {
			return err
		}
	}
	nfo.Message = pxy.Message
	nfo.Cause = cause
	nfo.Obj = pxy.Obj
	nfo.ObjKind = pxy.ObjKind
	return nil
}

// UnmarshalNotFound reads NotFound or it's subclasses from JSON bytes
func UnmarshalNotFound(in []byte) (NotFound, error) {
	fault, err := UnmarshalFault(in)
	if err != nil {
		return nil, err
	}
	if notFound, ok := fault.(NotFound); ok {
		return notFound, nil
	}
	return nil, fmt.Errorf("cannot unmarshal NotFound %v", fault)
}
