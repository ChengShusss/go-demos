package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

type OperationFunc func()

type Operation struct {
	Name string
	Desp string
	Func OperationFunc
}

type OperationSet struct {
	set map[string]*Operation
}

func (s *OperationSet) AddOperation(name, desp string, f OperationFunc) {
	s.set[name] = &Operation{name, desp, f}
}

func (s *OperationSet) ParseAndHandle(operation string) {
	op, ok := s.set[operation]
	if !ok {
		s.PrintInfo()
		return
	}
	op.Func()
}

var operationSet = OperationSet{
	set: map[string]*Operation{},
}

func (s *OperationSet) PrintInfo() {
	maxlen := 0
	for _, op := range s.set {
		if len(op.Name) > maxlen {
			maxlen = len(op.Name) + 5
		}
	}

	for _, op := range s.set {
		format := fmt.Sprintf("%%-%ds%%s\n", maxlen)
		fmt.Printf(format, op.Name, op.Desp)
	}

}

func main() {

	operationSet.AddOperation(
		"pack",
		"pack contents under specific folder into epub file",
		func() {})

	operationSet.AddOperation(
		"demo",
		"show the demos",
		demo)

	if len(os.Args) == 1 {
		operationSet.PrintInfo()
		return
	}

	operationSet.ParseAndHandle(os.Args[1])
}

func demo() {
	var all bool
	pflag.BoolVarP(&all, "all", "a", false, "add all files")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.CommandLine.Parse(os.Args[2:])
	fmt.Printf("all=%v\n", all)
}
