package models

import (
	"encoding/json"
	"fmt"

	"gitlab.eng.vmware.com/kkaraatanassov/go-json/interfaces"
)

// RuntimeFault represents fault
type RuntimeFault struct {
	Fault
}

var _ interfaces.Fault = &RuntimeFault{}
var _ interfaces.RuntimeFault = &RuntimeFault{}
var _ json.Marshaler = &RuntimeFault{}
var _ json.Unmarshaler = &RuntimeFault{}

// ZzRuntimeFault is a marker it prevents converting Fault struct to
// RuntimeFault interface
func (rf *RuntimeFault) ZzRuntimeFault() {
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
		Cause   FaultField
	}{}
	err := json.Unmarshal(in, pxy)
	if err != nil {
		return err
	}
	rf.Message = pxy.Message
	rf.Cause = pxy.Cause.Fault
	return nil
}

// RuntimeFaultField type allows reading polymorphic RuntimeFault fields
type RuntimeFaultField struct {
	interfaces.RuntimeFault
}

var _ interfaces.Fault = &RuntimeFaultField{}
var _ json.Unmarshaler = &RuntimeFaultField{}

// UnmarshalJSON reads the embedded fault taking care of the discriminator
func (ff *RuntimeFaultField) UnmarshalJSON(in []byte) error {
	var err error
	ff.RuntimeFault, err = UnmarshalRuntimeFault(in)
	return err
}

// UnmarshalRuntimeFault reads RuntimeFault or it's subclasses from JSON bytes
func UnmarshalRuntimeFault(in []byte) (interfaces.RuntimeFault, error) {
	fault, err := UnmarshalFault(in)
	if err != nil {
		return nil, err
	}
	if runtimeFault, ok := fault.(interfaces.RuntimeFault); ok {
		return runtimeFault, nil
	}
	return nil, fmt.Errorf("Cannot unmarshal RuntimeFault %v", fault)
}
