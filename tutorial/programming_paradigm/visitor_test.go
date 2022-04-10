package tutorial

import (
	"fmt"
	"testing"
)

//kubectl 使用了build模式和visitor模式
//https://github.com/kubernetes/kubernetes/blob/cea1d4e20b4a7886d8ff65f34c6d4f95efcb4742/staging/src/k8s.io/cli-runtime/pkg/resource/visitor.go
type VisitorFunc func(*Info, error) error

type Visitor interface {
	Visit(VisitorFunc) error
}

type Info struct {
	Namespace   string
	Name        string
	OtherThings string
}

func (info *Info) Visit(fn VisitorFunc) error {
	return fn(info, nil)
}

type NameVisitor struct {
	visitor Visitor
}

func (v NameVisitor) Visit(fn VisitorFunc) error {
	return v.visitor.Visit(func(info *Info, err error) error {
		fmt.Println("NameVisitor() before call function")
		err = fn(info, err)
		if err == nil {
			fmt.Printf("==> Name=%s, NameSpace=%s\n", info.Name, info.Namespace)
		}
		fmt.Println("NameVisitor() after call function")
		return err
	})
}

func WithNameVisitor() VisitorFunc {
	return func(info *Info, err error) error {
		fmt.Println("NameVisitor() before call function")
		info.Name = "kay.name"
		return nil
	}
}

type OtherVisitor struct {
	visitor Visitor
}

func (v OtherVisitor) Visit(fn VisitorFunc) error {
	return v.visitor.Visit(func(info *Info, err error) error {
		fmt.Println("OtherThingsVisitor() before call function")
		err = fn(info, err)
		if err == nil {
			fmt.Printf("==> OtherThings=%s\n", info.OtherThings)
		}
		fmt.Println("OtherThingsVisitor() after call function")
		return err
	})
}

func WithOtherVisitor() VisitorFunc {
	return func(info *Info, err error) error {
		fmt.Println("WithOtherVisitor() before call function")
		info.Namespace = "nanjing"
		return nil
	}
}

type DecoratedVisitor struct {
	visitor    Visitor
	decorators []VisitorFunc
}

func NewDecoratedVisitor(v Visitor, fn ...VisitorFunc) Visitor {
	if len(fn) == 0 {
		return v
	}
	return DecoratedVisitor{v, fn}
}

// Visit implements Visitor
func (v DecoratedVisitor) Visit(fn VisitorFunc) error {
	return v.visitor.Visit(func(info *Info, err error) error {
		if err != nil {
			return err
		}
		for i := range v.decorators {
			if err := v.decorators[i](info, nil); err != nil {
				return err
			}
		}
		if err := fn(info, nil); err != nil {
			return err
		}
		return nil
	})
}

func TestDecoratedVisitor(t *testing.T) {

	info := Info{}
	var v Visitor = &info

	v = NewDecoratedVisitor(v, WithNameVisitor(), WithOtherVisitor())
	loadFile := func(info *Info, err error) error {
		fmt.Println("process loadFile")
		info.OtherThings = "We are running as remote team."
		fmt.Printf("info: %v", info)
		return nil
	}
	v.Visit(loadFile)
}
