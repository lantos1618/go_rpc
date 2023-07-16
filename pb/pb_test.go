package pb_test

import (
	"bytes"
	"strings"
	"testing"

	"zug.dev/go_rpc/pb"
)

type SimpleStruct struct {
	Name string
	Age  int32
}

type CompoundStruct struct {
	Person SimpleStruct
	Phones []string
}

func MyFunc(name string, age int32) (string, error) {
	return "", nil
}

func SimpleFunc(person SimpleStruct) (string, error) {
	return "", nil
}

func CompoundFunc(person CompoundStruct) (string, error) {
	return "", nil
}

func TestMyFuncProtoGeneration(t *testing.T) {
	p := pb.NewProtoContext()
	v1 := p.Service("v1")
	v1.Register(MyFunc, "MyFunc")

	var b bytes.Buffer
	err := p.Dump(&b)
	if err != nil {
		t.Fatalf("Failed to generate proto file: %v", err)
	}

	expected := `syntax = "proto3";
service v1 {
  rpc MyFunc(MyFuncRequest) returns (MyFuncResponse);
}
// this is auto genorated for go's error struct
message Error {
  string Message = 1;
}
message MyFuncRequest {
  string arg1 = 1;
  int32 arg2 = 2;
}
message MyFuncResponse {
  string arg1 = 1;
  Error arg2 = 2;
}
`
	if strings.TrimSpace(b.String()) != strings.TrimSpace(expected) {
		t.Fatalf("Unexpected output:\n%s\nExpected:\n%s", b.String(), expected)
	}
}

// func TestSimpleFuncProtoGeneration(t *testing.T) {
// 	p := pb.NewProtoContext()
// 	v1 := p.Service("v1")
// 	v1.Register(SimpleFunc, "SimpleFunc")

// 	var b bytes.Buffer
// 	err := p.Dump(&b)
// 	if err != nil {
// 		t.Fatalf("Failed to generate proto file: %v", err)
// 	}

// 	expected := `syntax = "proto3";
// service v1 {
//   rpc SimpleFunc(SimpleFuncRequest) returns (SimpleFuncResponse);
// }
// // this is auto genorated for go's error struct
// message Error {
//   string Message = 1;
// }
// message SimpleStruct {
//   string Name = 1;
//   int32 Age = 2;
// }
// message SimpleFuncRequest {
//   SimpleStruct arg1 = 1;
// }
// message SimpleFuncResponse {
//   string arg1 = 1;
//   Error arg2 = 2;
// }
// `

// 	if strings.TrimSpace(b.String()) != strings.TrimSpace(expected) {
// 		t.Fatalf("Unexpected output:\n%s\nExpected:\n%s", b.String(), expected)
// 	}
// }

// func TestCompoundFuncProtoGeneration(t *testing.T) {
// 	p := pb.NewProtoContext()
// 	v1 := p.Service("v1")
// 	v1.Register(CompoundFunc, "CompoundFunc")

// 	var b bytes.Buffer
// 	err := p.Dump(&b)
// 	if err != nil {
// 		t.Fatalf("Failed to generate proto file: %v", err)
// 	}

// 	expected := `
// syntax = "proto3";
// service v1 {
// 	rpc CompoundFunc(CompoundFuncRequest) returns (CompoundFuncResponse);
// }
// // this is auto genorated for go's error struct
// message Error {
// 	string Message = 1;
// }
// message CompoundStruct {
// 	SimpleStruct Person = 1;
// 	repeated string Phones = 2;
// }
// message CompoundFuncRequest {
// 	CompoundStruct arg1 = 1;
// }
// message CompoundFuncResponse {
// 	string arg1 = 1;
// 	Error arg2 = 2;
// }
// `
// 	if strings.TrimSpace(b.String()) != strings.TrimSpace(expected) {
// 		t.Fatalf("Unexpected output:\n%s\nExpected:\n%s", b.String(), expected)
// 	}
// }
