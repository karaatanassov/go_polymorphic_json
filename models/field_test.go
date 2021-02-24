package models

import (
	"encoding/json"
	"testing"

	"gitlab.eng.vmware.com/kkaraatanassov/go-json/interfaces"
)

type Container struct {
	FaultField interfaces.Fault
}

var _ json.Unmarshaler = &Container{}

func (c *Container) UnmarshalJSON(in []byte) error {
	// Deserialize into temp object of utility class
	temp := struct {
		FaultField FaultField
	}{}
	err := json.Unmarshal(in, &temp)
	if err != nil {
		return err
	}
	// Re-assign all fields to the container unwrapping the util classes
	c.FaultField = temp.FaultField.Fault
	return nil
}

var container = Container{
	FaultField: notFound,
}

func TestContainerField(t *testing.T) {
	b, err := json.Marshal(&container)
	if err != nil {
		t.Error("Serialization failed", err)
	}

	s := string(b)
	t.Log("JSON Bytes", s)

	c := Container{}
	err = json.Unmarshal(b, &c)
	if err != nil {
		t.Error("Cannot deserialize fault", err)
		return
	}

	validateNotFound(c.FaultField, t)
}
