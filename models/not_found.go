package models

import (
	"encoding/json"

	"gitlab.eng.vmware.com/kkaraatanassov/go-json/interfaces"
)

// NotFound contains the data about a not found error
type NotFound struct {
	RuntimeFault
	ObjKind string
	Obj     string
}

var _ interfaces.NotFound = &NotFound{}
var _ interfaces.RuntimeFault = &NotFound{}
var _ interfaces.Fault = &NotFound{}
var _ json.Marshaler = &NotFound{}
var _ json.Unmarshaler = &NotFound{}

// GetObjKind retrieves the object kind of obj identifier
func (nfo *NotFound) GetObjKind() string {
	return nfo.ObjKind
}

// SetObjKind sets the kind of object references by obj
func (nfo *NotFound) SetObjKind(objKind string) {
	nfo.ObjKind = objKind
}

// GetObj retrieves the obj value
func (nfo *NotFound) GetObj() string {
	return nfo.Obj
}

// SetObj sets the obj id value
func (nfo *NotFound) SetObj(obj string) {
	nfo.Obj = obj
}

// MarshalJSON writes a NotFoundObject as JSON
func (nfo *NotFound) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Message string
		Cause   interfaces.Fault
		ObjKind string
		Obj     string
		Kind    string
	}{
		Message: nfo.Message,
		Cause:   nfo.Cause,
		ObjKind: nfo.ObjKind,
		Obj:     nfo.Obj,
		Kind:    "NotFound",
	})
}

// UnmarshalJSON reads a fault from JSON
func (nfo *NotFound) UnmarshalJSON(in []byte) error {
	pxy := &struct {
		Message string
		ObjKind string
		Obj     string
		Cause   FaultField
	}{}
	err := json.Unmarshal(in, pxy)
	if err != nil {
		return err
	}
	nfo.Message = pxy.Message
	nfo.Cause = pxy.Cause.Fault
	nfo.Obj = pxy.Obj
	nfo.ObjKind = pxy.ObjKind
	return nil
}
