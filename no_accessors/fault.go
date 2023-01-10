// Code to be generated.
// This is the only fully implemented binding type and includes:
// 1. Data object - Fault
// 2. Utility serialization object - FaultField
// 3. Array conversion utilities to be used in other bindings

package no_accessors

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// BaseFault is implemented by Error struct and included in RuntimeFault and
// NotFound. Thus one can upcast.
// This should be generated code.
type BaseFault interface {
	GetFault() *Fault
}

// Fault contains information about a base fault
// To be generated. This is actual data class. It implements the interface.Fault
// and the JSONSerializable
type Fault struct {
	Message string
	Cause   BaseFault
}

func init() {
	t["Fault"] = reflect.TypeOf((*Fault)(nil)).Elem()
}

var _ BaseFault = &Fault{}
var _ json.Marshaler = &Fault{}
var _ json.Unmarshaler = &Fault{}

// These assignments are not allowed
//var _ RuntimeFault = &Fault{}
//var _ NotFound = &Fault{}

func (f *Fault) GetFault() *Fault {
	return f
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
	var cause BaseFault
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
func UnmarshalFault(in []byte) (BaseFault, error) {
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

	reflectType, ok := t[kind]
	if !ok {
		return nil, fmt.Errorf("unknown type %v", kind)
	}

	var res BaseFault
	tv := reflect.New(reflectType)
	if res, ok = tv.Interface().(BaseFault); !ok {
		panic("A type in the registry does not implement Fault")
	}

	err = json.Unmarshal(in, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
