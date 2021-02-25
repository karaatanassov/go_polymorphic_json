package models

import (
	"encoding/json"
	"testing"

	"github.com/karaatanassov/go_polymorphic_json/interfaces"
)

type ArrayContainer struct {
	Faults []interfaces.Fault
}

var _ json.Unmarshaler = &ArrayContainer{}

func (c *ArrayContainer) UnmarshalJSON(in []byte) error {
	// Deserialize into temp object of utility class
	temp := struct {
		Faults []FaultField
	}{}
	err := json.Unmarshal(in, &temp)
	if err != nil {
		return err
	}
	// Re-assign all fields to the container unwrapping the util classes
	c.Faults = ToFaultsArray(temp.Faults)
	return nil
}

var arrayContainer ArrayContainer = ArrayContainer{
	Faults: []interfaces.Fault{
		fault,
		runtimeFault,
		notFound,
	},
}

func TestArray(t *testing.T) {
	b, err := json.Marshal(&arrayContainer)
	if err != nil {
		t.Error("Array serialization failed", err)
	}

	s := string(b)
	t.Log("JSON Bytes", s)

	c := ArrayContainer{}
	err = json.Unmarshal(b, &c)
	if err != nil {
		t.Error("Cannot deserialize fault array", err)
		return
	}
	if c.Faults == nil {
		t.Error("No faults were deserialized")
		return
	}
	if len(c.Faults) != 3 {
		t.Error("Expected 3 faults but encountered", len(c.Faults), c.Faults)
		return
	}
	validateFault(c.Faults[0], t)
	validateRuntimeFault(c.Faults[1], t)
	validateNotFound(c.Faults[2], t)
}
