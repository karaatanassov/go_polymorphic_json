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

// MarshalJSON writes a RuntimeFaultObject as JSON
func (rf *RuntimeFault) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Message string
		Kind    string
	}{
		Message: rf.Message,
		Kind:    "RuntimeFault",
	})
}
