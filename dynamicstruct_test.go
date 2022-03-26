package restflix

import (
	"encoding/json"
	"reflect"
	"testing"
)

type Container struct {
	Name  string
	Items []Item
}

type Item struct {
	Name      string
	Container Container
	Vals      []string
}

func TestDynamicStruct(t *testing.T) {
	//item := Item{}
	//c := reflect.Value{}
	//c.SetBool(true)
	//c.SetString("xd")
	//
	////val := reflect.ValueOf(item)
	//bx, _ := json.Marshal(c.Interface())
	//
	//t.Log(string(bx))
	//
	//return
	v0 := reflect.StructOf([]reflect.StructField{
		{
			Name: "C",
			Type: reflect.TypeOf(int(0)),
			Tag:  `json:"c"`,
		},
		{
			Name: "D",
			Type: reflect.TypeOf(""),
			Tag:  `json:"d"`,
		},
	})

	v := reflect.StructOf([]reflect.StructField{
		{
			Name: "A",
			Type: reflect.TypeOf(int(0)),
			Tag:  `json:"a"`,
		},
		{
			Name: "B",
			Type: reflect.TypeOf(""),
			Tag:  `json:"B"`,
		},
		{
			Name: "V0",
			Type: reflect.TypeOf(reflect.ValueOf(v0).Interface()),
			Tag:  `json:"v0"`,
		},
	})
	d := reflect.New(v).Interface()
	b0, _ := json.Marshal(d)

	t.Log(string(b0), reflect.ValueOf(v0).Interface())
	//
	//b0, _ := json.Marshal(d)
	//
	//instance1 := dynamicstruct.NewStruct().
	//	AddField("Z", "", `json:"z"`).
	//	AddField("E", 0.0, `json:"e"`).
	//	Build().
	//	New()
	//
	//b1, _ := json.Marshal(instance1)
	//
	//instance := dynamicstruct.NewStruct().
	//	AddField("X", "", `json:"x"`).
	//	AddField("Y", 0.0, `json:"y"`).
	//	AddField("Instance1", &b1, `json:"instance1"`).
	//	Build().
	//	New()
	//
	//b, _ := json.Marshal(instance)
	//
	//t.Log(string(b), string(b1), string(b0))
}
