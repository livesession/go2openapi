package ast

import "github.com/zdunecki/restflix/test/ast/objects"

// TODO: check pointers

//type MultiplePrimitiveTypes struct {
//	//String      string   `json:"string"`
//	//Boolean     bool     `json:"boolean"`
//	//Int         int      `json:"key"`
//	//Int32       int32    `json:"int32"`
//	//Float       float64  `json:"float"`
//	SliceString   []string       `json:"slice_string"`
//	SliceInt      []int          `json:"slice_int"`
//	SingleStructs []SingleStruct `json:"single_structs"`
//}

//type MultipleTypesOfSlices struct {
//	SliceString   []string       `json:"slice_string"`
//	SliceInt      []int          `json:"slice_int"`
//	SingleStructs []SingleStruct `json:"single_structs"`
//}

type TypeFromOutsideFileWithinPackage struct {
	//Response  OutsideFileResponse         `json:"response"`
	Response2 objects.OutsideFileResponse `json:"response2"`
}
