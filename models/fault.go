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
}

var _ interfaces.Fault = &Fault{}
var _ interfaces.JSONSerializable = &Fault{}

// GetMessage retrieves the message value
func (fault *Fault) GetMessage() string {
	return fault.Message
}

// SetMessage updates the message value
func (fault *Fault) SetMessage(message string) {
	fault.Message = message
}

// SerializeJSON writes Fault as JSON and adds discriminator
func (fault *Fault) SerializeJSON() ([]byte, error) {
	// The approach below passes the Fault ot the encoding/json code.
	// Thus this cannot be implemented in json.Marshaler.MashalJSON.
	// One alternative could be to copy the full object into a temporary object
	// with discriminator
	return json.Marshal(struct {
		*Fault
		Kind string
	}{
		Fault: fault,
		Kind:  "Fault",
	})
}

// DeserializeFault reads a fault from JSON and instantiates the proper type
// based on the Kind field. It deserializes the value twice. First scan for
// discriminator and then deserializes into the proper type.
func DeserializeFault(in []byte) (interfaces.Fault, error) {
	d := struct {
		Kind string
	}{}
	err := json.Unmarshal(in, &d)
	if err != nil {
		return nil, err
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

// FaultField is utility class that helps the go JSON deserializer to invoke the
// proper de-serialization logic for Fault fields while preserving the
// polymorphic nature of the type. In bindings we need two types - one with *Field
// and one for the actual object. UnmarshalJSON() reads into the the type with
// *Field fields and then copies the data. See field_test.go
type FaultField struct {
	interfaces.Fault
}

var _ interfaces.Fault = &FaultField{}
var _ json.Marshaler = &FaultField{}
var _ json.Unmarshaler = &FaultField{}

// MarshalJSON Writes the embedded Fault with discriminator
func (ff *FaultField) MarshalJSON() ([]byte, error) {
	return ff.Fault.SerializeJSON()
}

// UnmarshalJSON reads the embedded fault taking care of the discriminator
func (ff *FaultField) UnmarshalJSON(in []byte) error {
	var err error
	ff.Fault, err = DeserializeFault(in)
	return err
}

// ToFaultFieldArray is utlity to convert binding faults to serialization
// utility faults
func ToFaultFieldArray(faults []interfaces.Fault) []FaultField {
	var items []FaultField
	for _, tmp := range faults {
		items = append(items, FaultField{Fault: tmp})
	}
	return items
}

// ToFaultsArray is utlity to convert FaultField Array to interfaces.Fault array
func ToFaultsArray(faults []FaultField) []interfaces.Fault {
	var items []interfaces.Fault
	for _, tmp := range faults {
		items = append(items, tmp.Fault)
	}
	return items
}
