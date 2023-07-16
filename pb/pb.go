package pb

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

type namedType struct {
	name string
	typ  reflect.Type
}

type Service struct {
	name     string
	funcs    []namedType
	messages map[string]string
}

type ProtoContext struct {
	services map[string]*Service
}

func NewProtoContext() *ProtoContext {
	return &ProtoContext{
		services: make(map[string]*Service),
	}
}

func (p *ProtoContext) Service(name string) *Service {
	s, ok := p.services[name]
	if !ok {
		s = &Service{
			name:     name,
			funcs:    make([]namedType, 0),
			messages: make(map[string]string),
		}
		p.services[name] = s
		p.services[name].messages["Error"] = "// this is auto genorated for go's error struct\nmessage Error {\n  string Message = 1;\n}"

	}
	return s
}

func (s *Service) Register(fn interface{}, name string) {
	t := reflect.TypeOf(fn)
	if t.Kind() != reflect.Func {
		panic("register: not a func")
	}
	s.funcs = append(s.funcs, namedType{name: name, typ: t})

	// Generate input and output messages
	s.messages[fmt.Sprintf("%sRequest", name)] = s.genMessage(t, true, name)
	s.messages[fmt.Sprintf("%sResponse", name)] = s.genMessage(t, false, name)
}

func (s *Service) Gen() {
	fmt.Printf("service %s {\n", s.name)
	for _, nt := range s.funcs {
		fmt.Printf("  rpc %s(%sRequest) returns (%sResponse);\n", nt.name, nt.name, nt.name)
	}
	fmt.Println("}")
	for _, message := range s.messages {
		fmt.Println(message)
	}
}

func (s *Service) goTypeToProtoType(goType reflect.Type) string {
	switch goType.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int32:
		return "int32"
	case reflect.Int64:
		return "int64"
	case reflect.Slice:
		elemType := s.goTypeToProtoType(goType.Elem())
		return fmt.Sprintf("repeated %s", elemType)
	case reflect.Map:
		keyType := s.goTypeToProtoType(goType.Key())
		valueType := s.goTypeToProtoType(goType.Elem())
		return fmt.Sprintf("map<%s, %s>", keyType, valueType)
	case reflect.Struct:
		return goType.Name()
	case reflect.Interface:
		if goType.Name() == "error" {
			return "Error"
		}
		return "Unknown"
	default:
		return "Unknown"
	}
}
func (s *Service) genMessage(fnType reflect.Type, isInput bool, funcName string) string {
	var sb strings.Builder
	var start, end int
	var fieldType reflect.Type

	if isInput {
		sb.WriteString("message ")
		sb.WriteString(funcName + "Request")
		start, end = 0, fnType.NumIn()
	} else {
		sb.WriteString("message ")
		sb.WriteString(funcName + "Response")
		start, end = 0, fnType.NumOut()
	}

	sb.WriteString(" {\n")

	// Handle case of no input parameters or return values
	if end == 0 {
		sb.WriteString("}")
		return sb.String()
	}

	for i := start; i < end; i++ {
		if isInput {
			fieldType = fnType.In(i)
		} else {
			fieldType = fnType.Out(i)
		}
		if fieldType.Kind() == reflect.Struct {
			// If the field is a struct, generate a message for it
			if _, ok := s.messages[fieldType.Name()]; !ok {
				s.messages[fieldType.Name()] = s.genMessageStruct(fieldType)
			}
			sb.WriteString(fmt.Sprintf("  %s arg%d = %d;\n", fieldType.Name(), i+1, i+1))
		} else {
			sb.WriteString(fmt.Sprintf("  %s arg%d = %d;\n", s.goTypeToProtoType(fieldType), i+1, i+1))
		}
	}
	sb.WriteString("}")
	return sb.String()
}
func (s *Service) genMessageStruct(goType reflect.Type) string {
	var sb strings.Builder
	sb.WriteString("message ")
	sb.WriteString(goType.Name())
	sb.WriteString(" {\n")
	for i := 0; i < goType.NumField(); i++ {
		field := goType.Field(i)
		fieldType := s.goTypeToProtoType(field.Type)
		sb.WriteString(fmt.Sprintf("  %s %s = %d;\n", fieldType, field.Name, i+1))
	}
	sb.WriteString("}")
	return sb.String()
}

func (p *ProtoContext) Dump(w io.Writer) error {
	// Specify that we're using proto3
	if _, err := fmt.Fprintln(w, `syntax = "proto3";`); err != nil {
		return err
	}
	// Generate service and message definitions
	for _, service := range p.services {
		if err := service.Dump(w); err != nil {
			return err
		}
	}
	return nil
}
func (s *Service) Dump(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "service %s {\n", s.name); err != nil {
		return err
	}
	for _, nt := range s.funcs {
		if _, err := fmt.Fprintf(w, "  rpc %s(%sRequest) returns (%sResponse);\n", nt.name, nt.name, nt.name); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(w, "}"); err != nil {
		return err
	}
	for _, message := range s.messages {
		if _, err := fmt.Fprintln(w, message); err != nil {
			return err
		}
	}
	return nil
}

func (p *ProtoContext) DumpToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return p.Dump(file)
}
