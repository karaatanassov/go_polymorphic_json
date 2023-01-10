package no_accessors

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// BaseNotFound represents error when object is not found
// To be generated
type BaseNotFound interface {
	BaseRuntimeFault
	GetNotFound() *NotFound
}

// NotFound contains the data about a not found error
type NotFound struct {
	RuntimeFault
	ObjKind string
	Obj     string
}

func init() {
	t["NotFound"] = reflect.TypeOf((*NotFound)(nil)).Elem()
}

var _ BaseNotFound = &NotFound{}
var _ BaseRuntimeFault = &NotFound{}
var _ BaseFault = &NotFound{}
var _ json.Marshaler = &NotFound{}
var _ json.Unmarshaler = &NotFound{}

func (f *NotFound) GetNotFound() *NotFound {
	return f
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
		Cause   json.RawMessage
		ObjKind string
		Obj     string
	}{}
	err := json.Unmarshal(in, pxy)
	if err != nil {
		return err
	}
	var cause BaseFault
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
func UnmarshalNotFound(in []byte) (BaseNotFound, error) {
	fault, err := UnmarshalFault(in)
	if err != nil {
		return nil, err
	}
	if notFound, ok := fault.(BaseNotFound); ok {
		return notFound, nil
	}
	return nil, fmt.Errorf("cannot unmarshal NotFound %v", fault)
}
