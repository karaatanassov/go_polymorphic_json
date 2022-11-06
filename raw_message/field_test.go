package raw_message

import (
	"encoding/json"
	"testing"
)

type Container struct {
	FaultField Fault
}

var _ json.Unmarshaler = &Container{}

func (c *Container) UnmarshalJSON(in []byte) error {
	// Deserialize into temp object
	temp := struct {
		FaultField json.RawMessage
	}{}
	err := json.Unmarshal(in, &temp)
	if err != nil {
		return err
	}
	c.FaultField = nil
	if temp.FaultField != nil {
		c.FaultField, err = UnmarshalFault(temp.FaultField)
		if err != nil {
			return err
		}
	}
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
