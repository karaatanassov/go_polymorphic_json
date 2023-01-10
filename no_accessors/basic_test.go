package no_accessors

import (
	"encoding/json"
	"testing"
)

var innerFault = Fault{
	Message: "inner message",
}

var innerRuntimeFault = RuntimeFault{
	Fault: innerFault,
}

var fault *Fault = &Fault{
	Message: "test message",
	Cause:   &innerRuntimeFault,
}
var runtimeFault *RuntimeFault = &RuntimeFault{
	Fault: *fault,
}
var notFound *NotFound = &NotFound{
	RuntimeFault: *runtimeFault,
	ObjKind:      "VirtualMachine",
	Obj:          "vm-42",
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

func serializeDeserialize(s interface{}, t *testing.T) (BaseFault, error) {
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

func validateFault(fault BaseFault, t *testing.T) {
	switch v := fault.(type) {
	case *Fault:
		if v.Message != "test message" {
			t.Error("Unexpected message:", v.Message)
		}
		validateCause(v.Cause, t)
	default:
		t.Error("Unexpected type")
	}
}

func validateRuntimeFault(fault BaseFault, t *testing.T) {
	switch v := fault.(type) {
	case *RuntimeFault:
		if v.Message != "test message" {
			t.Error("Unexpected message:", v.Message)
		}
		validateCause(v.Cause, t)
	default:
		t.Error("Unexpected type")
	}
}

func validateNotFound(fault BaseFault, t *testing.T) {
	if notFound, ok := fault.(BaseNotFound); ok {
		nf := notFound.GetNotFound()
		if nf.Message != "test message" {
			t.Error("Unexpected message:", nf.Message)
		}
		if nf.ObjKind != "VirtualMachine" {
			t.Error("Unexpected obj kind:", nf.ObjKind)
		}
		if nf.Obj != "vm-42" {
			t.Error("Unexpected obj:", nf.Obj)
		}
		validateCause(nf.Cause, t)
	} else {
		t.Error("Unexpected type", fault)

	}
}

func validateCause(cause BaseFault, t *testing.T) {
	if cause == nil {
		t.Error("Missing cause")
	} else {
		switch i := cause.(type) {
		case *RuntimeFault:
			if i.Message != "inner message" {
				t.Error("Inner message is wrong:", i.Message)
			}
		default:
			t.Error("Unexpected inner type", cause)
		}
	}
}
