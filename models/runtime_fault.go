package models

import (
	"encoding/json"
	"fmt"

	"github.com/karaatanassov/go_polymorphic_json/interfaces"
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
	rf.Message = pxy.Message
	rf.Cause = cause
	return nil
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
