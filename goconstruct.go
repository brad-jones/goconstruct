package goconstruct

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/structs"
	"github.com/iancoleman/orderedmap"
	"github.com/wesovilabs/koazee"
)

// Construct is a class like object, create new instances with `New`.
type Construct struct {
	ID       string
	Type     string
	Props    interface{}
	Outputs  interface{}
	Parent   *Construct
	Children []*Construct
}

// New is the constructor for Construct.
//
// It uses a functional options based API, see:
// https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
//
// Example Usage:
// 	goconstruct.New("MyThing", goconstruct.Scope(parentThing),
// 		goconstruct.Props(&MyThingProps{}),
// 	)
func New(id string, options ...func(*Construct)) *Construct {
	this := &Construct{ID: id}
	for _, option := range options {
		option(this)
	}
	return this
}

// Scope defines the parent child relationship between constructs
func Scope(parent *Construct) func(*Construct) {
	return func(c *Construct) {
		parent.AddChild(c)
	}
}

// Props allows you to assign any arbitrary object, usually a struct
// specifically defined for your construct. So long as it can be marshalled
// to JSON that is all that matters.
func Props(v interface{}) func(*Construct) {
	return func(c *Construct) {
		c.Props = v
	}
}

// Outputs will take a pointer to a struct and ensure it is initialized with
// reference strings (a bit like a JSON Pointer). These strings will later be
// used by a generator to connect resources together in the appropriate order.
func Outputs(v interface{}) func(*Construct) {
	return func(c *Construct) {
		for _, f := range structs.Fields(v) {
			if !f.IsExported() {
				continue
			}
			if f.IsEmbedded() {
				if err := f.Set(c.Props); err != nil {
					panic(fmt.Errorf("%w", err))
				}
			} else {
				if f.Kind().String() == "string" {
					if err := f.Set(fmt.Sprintf("%s.%s", c.ID, f.Name())); err != nil {
						panic(fmt.Errorf("%w", err))
					}
				}
			}
		}
		c.Outputs = v
	}
}

// Type sets the Type field of a Construct
func Type(t string) func(*Construct) {
	return func(c *Construct) {
		if c.Type == "" {
			c.Type = t
		} else {
			c.Type = fmt.Sprintf("%s/%s", c.Type, t)
		}
	}
}

// Constructor expects a function that will be called with the Construct
// instance as it's only parameter in order to add new children to it.
func Constructor(builder func(*Construct)) func(*Construct) {
	return func(c *Construct) {
		builder(c)
	}
}

// ConstructorWithOutputs expects a function that will be called with the
// Construct instance as it's only parameter in order to add new children to it.
// It also expects a value to be returned with will be assigned to the Outputs.
func ConstructorWithOutputs(builder func(*Construct) interface{}) func(*Construct) {
	return func(c *Construct) {
		c.Outputs = builder(c)
	}
}

// AddChild will create the parent and child relationship between all Constructs
func (c *Construct) AddChild(child *Construct) {
	child.Parent = c
	if c.Children == nil {
		c.Children = []*Construct{}
	}
	c.Children = append(c.Children, child)
	child.ID = strings.Replace(
		fmt.Sprintf("%s/%s", c.ID, child.ID),
		"//", "/", 1,
	)
}

// Synth will convert the Construct graph into a serializable map
func (c *Construct) Synth() *orderedmap.OrderedMap {
	o := orderedmap.New()
	o.Set("id", c.ID)
	o.Set("type", c.Type)
	o.Set("props", c.Props)
	o.Set("children", koazee.StreamOf(c.Children).
		Map(func(v *Construct) *orderedmap.OrderedMap { return v.Synth() }).
		Do().Out().Val(),
	)
	return o
}

// MarshalJSON implements the Marshaler interface
// see https://golang.org/pkg/encoding/json/#Marshaler
func (c *Construct) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Synth())
}
