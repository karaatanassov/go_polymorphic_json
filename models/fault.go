// Code to be generated for a base type.

// Package models contains the various data structures
package models

import (
	"encoding/json"

	"github.com/karaatanassov/go_polymorphic_json/interfaces"
)

// Fault contains information about a base fault
// To be generated. This is actual data class. It implements the interface.Fault
// and the JSONSerializable
type Fault struct {
	Message string
	Cause   interfaces.Fault
}

var _ interfaces.Fault = &Fault{}
var _ json.Marshaler = &Fault{}
var _ json.Unmarshaler = &Fault{}

// This assignment is not allowed with ZzRuntimeFault
//var _ interfaces.RuntimeFault = &Fault{}

// This assignment is not allowed
//var _ interfaces.NotFound = &Fault{}

// GetMessage retrieves the message value
func (fault *Fault) GetMessage() string {
	return fault.Message
}

// SetMessage updates the message value
func (fault *Fault) SetMessage(message string) {
	fault.Message = message
}

// GetCause returns the case of fault
func (fault *Fault) GetCause() interfaces.Fault {
	return fault.Cause
}

// SetCause sets the cause of the fault
func (fault *Fault) SetCause(cause interfaces.Fault) {
	fault.Cause = cause
}

// UnmarshalJSON reads a fault from JSON
func (fault *Fault) UnmarshalJSON(in []byte) error {
	pxy := &struct {
		Message string
		Cause   json.RawMessage
	}{}
	err := json.Unmarshal(in, pxy)
	if err != nil {
		return err
	}
	var cause interfaces.Fault
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
func (fault *Fault) MarshalJSON() ([]byte, error) {
	type marshalable Fault
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
func UnmarshalFault(in []byte) (interfaces.Fault, error) {
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

	var res interfaces.Fault
	switch kind {
	case "NotFound":
		res = &NotFound{}
	case "RuntimeFault":
		res = &RuntimeFault{}
	default: // Error on default or try to use base type?
		res = &Fault{}
	}
	json.Unmarshal(in, res)

	return res, nil
}
