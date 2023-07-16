package main

import (
	"fmt"
	"log"

	"zug.dev/go_rpc/pb"
)

// Define the structs for the tests
type SimpleStruct struct {
	Name string
	Age  int32
}

type CompoundStruct struct {
	Person SimpleStruct
	Phones []string
}

func MyFunc(name string, age int32) (string, error) {
	return fmt.Sprintf("Hello %s, you are %d years old", name, age), nil
}

func SimpleFunc(person SimpleStruct) (string, error) {
	return fmt.Sprintf("Hello %s, you are %d years old", person.Name, person.Age), nil
}

func CompoundFunc(person CompoundStruct) (string, error) {
	return fmt.Sprintf("Hello %s, you are %d years old", person.Person.Name, person.Person.Age), nil
}

func main() {
	p := pb.NewProtoContext()
	v1 := p.Service("v1")
	v1.Register(MyFunc, "MyFunc")
	v1.Register(SimpleFunc, "SimpleFunc")
	v1.Register(CompoundFunc, "CompoundFunc")
	if err := p.DumpToFile("rpc.proto"); err != nil {
		log.Fatal(err)
	}
}
