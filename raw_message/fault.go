// Code to be generated for a base type.

// Package models contains the various data structures
package raw_message

import (
	"encoding/json"
)

// Fault represents a base error
// This should be generated code. It is possible to emit it in the models
// package and creative name will be needed e.g. NotFoundInterface.
// The interface package approach seems cleaner.
type Fault interface {
	GetMessage() string
	SetMessage(string)
	GetCause() Fault
	SetCause(Fault)
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

// This assignment is not allowed with ZzRuntimeFault
//var _ RuntimeFault = &Fault{}

// This assignment is not allowed
//var _ NotFound = &Fault{}

// GetMessage retrieves the message value
func (fault *FaultStruct) GetMessage() string {
	return fault.Message
}

// SetMessage updates the message value
func (fault *FaultStruct) SetMessage(message string) {
	fault.Message = message
}

// GetCause returns the case of fault
func (fault *FaultStruct) GetCause() Fault {
	return fault.Cause
}

// SetCause sets the cause of the fault
func (fault *FaultStruct) SetCause(cause Fault) {
	fault.Cause = cause
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
	json.Unmarshal(in, res)

	return res, nil
}
