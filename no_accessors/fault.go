// Code to be generated.
// This is the only fully implemented binding type and includes:
// 1. Data object - Fault
// 2. Utility serialization object - FaultField
// 3. Array conversion utilities to be used in other bindings

package no_accessors

import (
	"encoding/json"
)

// Fault is implemented by Error struct and included in RuntimeFault and
// NotFound. Thus one can upcast.
// This should be generated code.
type Fault interface {
	GetFault() *FaultStruct
}

// FaultStruct contains information about a base fault
// To be generated. This is actual data class. It implements the interface.FaultStruct
// and the JSONSerializable
type FaultStruct struct {
	Message string
	Cause   Fault
}

var _ Fault = &FaultStruct{}
var _ json.Marshaler = &FaultStruct{}
var _ json.Unmarshaler = &FaultStruct{}

// This assignment is not allowed
//var _ RuntimeFault = &Fault{}

// This assignment is not allowed
//var _ NotFound = &Fault{}

func (f *FaultStruct) GetFault() *FaultStruct {
	return f
}

// UnmarshalJSON reads a fault from JSON
func (fault *FaultStruct) UnmarshalJSON(in []byte) error {
	pxy := &struct {
		Message string
		Cause   json.RawMessage
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
	fault.Message = pxy.Message
	fault.Cause = cause
	return nil
}

// MarshalJSON writes Fault as JSON and adds discriminator
func (fault *FaultStruct) MarshalJSON() ([]byte, error) {
	type marshalable FaultStruct
	// The approach below copies the full object into a temporary object
	// with discriminator and passes it to the go json mashaler.
	// An alternative is to create small utility object that holds the
	// discriminator and anonymous Fault pointer. This requires the
	// serialization to be in custom interface and additional copy logic in
	// higher level bindings. The current approach preserves the go abstractions
	// and simplifies higher level bindings.
	return json.Marshal(struct {
		Kind string
		marshalable
	}{
		Kind:        "Fault",
		marshalable: marshalable(*fault),
	})
}

// UnmarshalFault reads a fault from JSON and instantiates the proper type
// based on the Kind field. It deserializes the value twice. First scan for
// discriminator and then deserializes into the proper type.
func UnmarshalFault(in []byte) (Fault, error) {
	d := &struct {
		Kind string
	}{}
	// Double pointer detects null values
	err := json.Unmarshal(in, &d)
	if err != nil {
		return nil, err
	}
	if d == nil {
		return nil, nil
	}
	kind := d.Kind

	var res Fault
	switch kind {
	case "NotFound":
		res = &NotFoundStruct{}
	case "RuntimeFault":
		res = &RuntimeFaultStruct{}
	default: // Error on default or try to use base type?
		res = &FaultStruct{}
	}
	err = json.Unmarshal(in, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
