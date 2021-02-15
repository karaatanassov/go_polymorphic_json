package models

import (
	"encoding/json"

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

// MarshalJSON writes a RuntimeFaultObject as JSON
func (rf *RuntimeFault) MarshalJSON() ([]byte, error) {
	type marshalable RuntimeFault
	return json.Marshal(struct {
		marshalable
		Kind string
	}{
		marshalable: marshalable(*rf),
		Kind:        "RuntimeFault",
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
