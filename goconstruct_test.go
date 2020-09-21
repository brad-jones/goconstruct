package goconstruct_test

import (
	"testing"

	"github.com/brad-jones/goconstruct"
	"github.com/stretchr/testify/assert"
)

func TestId(t *testing.T) {
	t.Parallel()
	c := goconstruct.New("Foo")
	assert.Equal(t, "Foo", c.ID)
}

func TestScope(t *testing.T) {
	t.Parallel()
	p := goconstruct.New("Foo")
	c := goconstruct.New("Bar", goconstruct.Scope(p))
	assert.Equal(t, p, c.Parent)
	assert.Equal(t, c, p.Children[0])
	assert.Equal(t, "Foo/Bar", c.ID)
}

type FooProps struct {
	Bar string
}

func TestProps(t *testing.T) {
	t.Parallel()
	c := goconstruct.New("Foo", goconstruct.Props(&FooProps{Bar: "baz"}))
	assert.Equal(t, "baz", c.Props.(*FooProps).Bar)
}

type FooOutputs struct {
	Bar string
}

func TestOutputs(t *testing.T) {
	t.Parallel()
	c := goconstruct.New("Foo", goconstruct.Outputs(&FooOutputs{}))
	assert.Equal(t, "Foo.Bar", c.Outputs.(*FooOutputs).Bar)
}

func TestType(t *testing.T) {
	t.Parallel()
	c := goconstruct.New("Foo", goconstruct.Type("Bar"))
	assert.Equal(t, "Bar", c.Type)
}

func TestTemplatedType(t *testing.T) {
	t.Parallel()
	c := goconstruct.New("Foo", goconstruct.Type("Bar"), goconstruct.Type("Baz"))
	assert.Equal(t, "Bar/Baz", c.Type)
}

func TestConstructor(t *testing.T) {
	t.Parallel()
	goconstruct.New("Foo", goconstruct.Constructor(func(c *goconstruct.Construct) {
		assert.Equal(t, "Foo", c.ID)
	}))
}

func TestConstructorWithOutputs(t *testing.T) {
	t.Parallel()
	c := goconstruct.New("Foo", goconstruct.ConstructorWithOutputs(func(c *goconstruct.Construct) interface{} {
		assert.Equal(t, "Foo", c.ID)
		return "Bar"
	}))
	assert.Equal(t, "Bar", c.Outputs.(string))
}

func TestMarshalJSON(t *testing.T) {
	t.Parallel()
	c := goconstruct.New("Foo", goconstruct.Constructor(func(c *goconstruct.Construct) {
		goconstruct.New("Bar", goconstruct.Scope(c), goconstruct.Constructor(func(c *goconstruct.Construct) {
			goconstruct.New("Baz", goconstruct.Scope(c))
		}))
	}))
	json, err := c.MarshalJSON()
	if assert.NoError(t, err) {
		assert.JSONEq(t, `{
			"id": "Foo",
			"type": "",
			"props": null,
			"children": [
				{
					"id": "Foo/Bar",
					"type": "",
					"props": null,
					"children": [
						{
							"id": "Foo/Bar/Baz",
							"type": "",
							"props": null,
							"children": []
						}
					]
				}
			]
		}`, string(json))
	}
}
