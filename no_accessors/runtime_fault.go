package no_accessors

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// BaseRuntimeFault represents all runtime faults that can be thrown
// To be generated
type BaseRuntimeFault interface {
	BaseFault
	GetRuntimeFault() *RuntimeFault
}

// RuntimeFault represents fault
type RuntimeFault struct {
	Fault
}

func init() {
	t["RuntimeFault"] = reflect.TypeOf((*RuntimeFault)(nil)).Elem()
}

var _ BaseFault = &RuntimeFault{}
var _ BaseRuntimeFault = &RuntimeFault{}
var _ json.Marshaler = &RuntimeFault{}
var _ json.Unmarshaler = &RuntimeFault{}

func (f *RuntimeFault) GetRuntimeFault() *RuntimeFault {
	return f
}

// MarshalJSON writes a RuntimeFaultObject as JSON
func (rf *RuntimeFault) MarshalJSON() ([]byte, error) {
	type marshalable RuntimeFault
	return json.Marshal(struct {
		Kind string
		marshalable
	}{
		Kind:        "RuntimeFault",
		marshalable: marshalable(*rf),
	})
}

// UnmarshalJSON reads a fault from JSON
func (rf *RuntimeFault) UnmarshalJSON(in []byte) error {
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
	rf.Message = pxy.Message
	rf.Cause = cause
	return nil
}

// UnmarshalRuntimeFault reads RuntimeFault or it's subclasses from JSON bytes
func UnmarshalRuntimeFault(in []byte) (BaseRuntimeFault, error) {
	fault, err := UnmarshalFault(in)
	if err != nil {
		return nil, err
	}
	if runtimeFault, ok := fault.(BaseRuntimeFault); ok {
		return runtimeFault, nil
	}
	return nil, fmt.Errorf("cannot unmarshal RuntimeFault %v", fault)
}
