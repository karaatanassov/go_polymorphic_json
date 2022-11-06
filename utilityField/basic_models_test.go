package models

import (
	"encoding/json"
	"testing"
)

var innerFault = FaultStruct{
	Message: "inner message",
}

var innerRuntimeFault = RuntimeFaultStruct{
	FaultStruct: innerFault,
}

var fault *FaultStruct = &FaultStruct{
	Message: "test message",
	Cause:   &innerRuntimeFault,
}
var runtimeFault *RuntimeFaultStruct = &RuntimeFaultStruct{
	FaultStruct: *fault,
}
var notFound *NotFoundStruct = &NotFoundStruct{
	RuntimeFaultStruct: *runtimeFault,
	ObjKind:            "VirtualMachine",
	Obj:                "vm-42",
}

func TestFault(t *testing.T) {
	fault, err := serializeDeserialize(fault, t)
	if err != nil {
		t.Error("Fault basic test failed", err)
		return
	}
	validateFault(fault, t)
}

func TestRuntimeFault(t *testing.T) {
	fault, err := serializeDeserialize(runtimeFault, t)
	if err != nil {
		t.Error("RuntimeFault basic test failed", err)
		return
	}
	validateRuntimeFault(fault, t)
}

func TestNotFound(t *testing.T) {
	fault, err := serializeDeserialize(notFound, t)
	if err != nil {
		t.Error("NotFound basic test failed", err)
		return
	}
	validateNotFound(fault, t)
}

func TestInvalidRuntimeFault(t *testing.T) {
	b, err := json.Marshal(fault)
	if err != nil {
		t.Error("Serialization failed", err)
		return
	}

	t.Log("JSON Bytes", string(b))

	_, err = UnmarshalRuntimeFault(b)
	if err == nil {
		t.Error("Expected to fail unmarshaling fault into RuntimeFault")
	}
}

func TestValidRuntimeFault(t *testing.T) {
	b, err := json.Marshal(notFound)
	if err != nil {
		t.Error("Serialization failed", err)
		return
	}

	t.Log("JSON Bytes", string(b))

	rtf, err := UnmarshalRuntimeFault(b)
	if err != nil {
		t.Error("Expected to unmarshal RuntimeFault.", err)
	}
	// Note that the validate method accepts Fault and works well with
	// RuntimeFault which holds NotFound
	validateNotFound(rtf, t)
}

func serializeDeserialize(s interface{}, t *testing.T) (Fault, error) {
	b, err := json.Marshal(s)
	if err != nil {
		t.Error("Serialization failed", err)
		return nil, err
	}

	t.Log("JSON Bytes", string(b))

	fault, err := UnmarshalFault(b)
	if err != nil {
		t.Error("Cannot deserialize fault", err)
		return nil, err
	}
	return fault, nil
}

func validateFault(fault Fault, t *testing.T) {
	switch v := fault.(type) {
	case *FaultStruct:
		if v.Message != "test message" {
			t.Error("Unexpected message:", v.Message)
		}
		validateCause(v.Cause, t)
	default:
		t.Error("Unexpected type")
	}
}

func validateRuntimeFault(fault Fault, t *testing.T) {
	switch v := fault.(type) {
	case *RuntimeFaultStruct:
		if v.Message != "test message" {
			t.Error("Unexpected message:", v.Message)
		}
		validateCause(v.Cause, t)
	default:
		t.Error("Unexpected type")
	}
}

func validateNotFound(fault Fault, t *testing.T) {
	if fault.GetMessage() != "test message" {
		t.Error("Unexpected message:", fault.GetMessage())
	}

	if notFound, ok := fault.(NotFound); ok {
		if notFound.GetObjKind() != "VirtualMachine" {
			t.Error("Unexpected obj kind:", notFound.GetObjKind())
		}
		if notFound.GetObj() != "vm-42" {
			t.Error("Unexpected obj:", notFound.GetObj())
		}
		validateCause(notFound.GetCause(), t)
	} else {
		t.Error("Unexpected type", fault)

	}
}

func validateCause(cause Fault, t *testing.T) {
	if cause == nil {
		t.Error("Missing cause")
	} else {
		switch i := cause.(type) {
		case *RuntimeFaultStruct:
			if i.Message != "inner message" {
				t.Error("Inner message is wrong:", i.Message)
			}
		default:
			t.Error("Unexpected inner type", cause)
		}
	}
}
