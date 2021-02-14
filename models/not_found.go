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
var _ interfaces.JSONSerializable = &NotFound{}

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

// SerializeJSON writes a NotFoundObject as JSON
func (nfo *NotFound) SerializeJSON() ([]byte, error) {
	return json.Marshal(struct {
		*NotFound
		Kind string
	}{
		NotFound: nfo,
		Kind:     "NotFound",
	})
}
