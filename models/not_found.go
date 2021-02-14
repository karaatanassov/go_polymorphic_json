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
		Message *string
		ObjKind *string
		Obj     *string
		Kind    string
	}{
		Message: &nfo.Message,
		ObjKind: &nfo.ObjKind,
		Obj:     &nfo.Obj,
		Kind:    "NotFound",
	})
}
