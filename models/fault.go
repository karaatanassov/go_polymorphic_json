// Code to be generated.
// This is the only fully implemented binding type and includes:
// 1. Data object - Fault
// 2. Utility serialization object - FaultField
// 3. Array conversion utilities to be used in other bindings

package models

import (
	"encoding/json"

	"gitlab.eng.vmware.com/kkaraatanassov/go-json/interfaces"
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
		Cause   FaultField
	}{}
	err := json.Unmarshal(in, pxy)
	if err != nil {
		return err
	}
	fault.Message = pxy.Message
	fault.Cause = pxy.Cause.Fault
	return nil
}

// MarshalJSON writes Fault as JSON and adds discriminator
func (fault *Fault) MarshalJSON() ([]byte, error) {
	type marshalFault Fault
	return json.Marshal(struct {
		marshalFault
		Kind string
	}{
		marshalFault: marshalFault(*fault),
		Kind:         "Fault",
	})
}

// DeserializeFault reads a fault from JSON and instantiates the proper type
// based on the Kind field. It deserializes the value twice. First scan for
// discriminator and then deserializes into the proper type.
func DeserializeFault(in []byte) (interfaces.Fault, error) {
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
	err = json.Unmarshal(in, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// FaultField is utility class that helps the go JSON deserializer to invoke the
// proper de-serialization logic for Fault fields while preserving the
// polymorphic nature of the type. go uses reflection to invoke the proper
// de-serialization method. As interface do not have methods implementations
// we need a concrete class field that will have the logic to deserialize the
// proper interface implementation.
// In bindings deserialization we need two types - one with *Field
// and one for the interface type. UnmarshalJSON() reads into the the * Field
// type and then copies the data to the type with interface type.
// See field_test.go
type FaultField struct {
	interfaces.Fault
}

var _ interfaces.Fault = &FaultField{}
var _ json.Unmarshaler = &FaultField{}

// UnmarshalJSON reads the embedded fault taking care of the discriminator
func (ff *FaultField) UnmarshalJSON(in []byte) error {
	var err error
	ff.Fault, err = DeserializeFault(in)
	return err
}

// ToFaultsArray is utility to convert FaultField Array to interfaces.Fault array
func ToFaultsArray(faults []FaultField) []interfaces.Fault {
	var items []interfaces.Fault
	for _, tmp := range faults {
		items = append(items, tmp.Fault)
	}
	return items
}
