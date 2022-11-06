package raw_message

import (
	"encoding/json"
	"fmt"
)

// RuntimeFault represents all runtime faults that can be thrown
// To be generated
type RuntimeFault interface {
	Fault
	// ZzRuntimeFault disallows converting Fault struct to RuntimeFault interface
	ZzRuntimeFault()
}

// RuntimeFaultStruct represents fault
type RuntimeFaultStruct struct {
	FaultStruct
}

var _ Fault = &RuntimeFaultStruct{}
var _ RuntimeFault = &RuntimeFaultStruct{}
var _ json.Marshaler = &RuntimeFaultStruct{}
var _ json.Unmarshaler = &RuntimeFaultStruct{}

// ZzRuntimeFault is a marker it prevents converting Fault struct to
// RuntimeFault interface
func (rf *RuntimeFaultStruct) ZzRuntimeFault() {
}

// MarshalJSON writes a RuntimeFaultObject as JSON
func (rf *RuntimeFaultStruct) MarshalJSON() ([]byte, error) {
	type marshalable RuntimeFaultStruct
	return json.Marshal(struct {
		Kind string
		marshalable
	}{
		Kind:        "RuntimeFault",
		marshalable: marshalable(*rf),
	})
}

// UnmarshalJSON reads a fault from JSON
func (rf *RuntimeFaultStruct) UnmarshalJSON(in []byte) error {
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
	rf.Message = pxy.Message
	rf.Cause = cause
	return nil
}

// UnmarshalRuntimeFault reads RuntimeFault or it's subclasses from JSON bytes
func UnmarshalRuntimeFault(in []byte) (RuntimeFault, error) {
	fault, err := UnmarshalFault(in)
	if err != nil {
		return nil, err
	}
	if runtimeFault, ok := fault.(RuntimeFault); ok {
		return runtimeFault, nil
	}
	return nil, fmt.Errorf("cannot unmarshal RuntimeFault %v", fault)
}
