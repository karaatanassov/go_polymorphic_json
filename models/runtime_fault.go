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
var _ interfaces.JSONSerializable = &RuntimeFault{}

// SerializeJSON writes a RuntimeFaultObject as JSON
func (rf *RuntimeFault) SerializeJSON() ([]byte, error) {
	return json.Marshal(struct {
		*RuntimeFault
		Kind string
	}{
		RuntimeFault: rf,
		Kind:         "RuntimeFault",
	})
}
