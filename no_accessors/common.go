package no_accessors

import "reflect"

var (
	// t is a global type map for unmarhsaling polymorphic types
	t map[string]reflect.Type
)

func init() {
	t = make(map[string]reflect.Type)
}
