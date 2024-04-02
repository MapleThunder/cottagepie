package object

import (
	"bytes"
	"cottagepie/ast"
	"fmt"
	"hash/fnv"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	STRING_OBJ       = "STRING"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	SERVES_VALUE_OBJ = "SERVES_VALUE"
	ERROR_OBJ        = "ERROR"
	RECIPE_OBJ       = "RECIPE"
	BUILT_IN_OBJ     = "BUILT_IN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) HashKey() HashKey { return HashKey{Type: i.Type(), Value: uint64(i.Value)} }

// String
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// Boolean
type Boolean struct {
	Value bool
}

func (i *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (i *Boolean) Inspect() string  { return fmt.Sprintf("%t", i.Value) }
func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

// Null
type Null struct{}

func (i *Null) Type() ObjectType { return NULL_OBJ }
func (i *Null) Inspect() string  { return "null" }

// Serves Value
type ServesValue struct {
	Value Object
}

func (sv *ServesValue) Type() ObjectType { return SERVES_VALUE_OBJ }
func (sv *ServesValue) Inspect() string  { return sv.Value.Inspect() }

// Error
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR>> " + e.Message }

// Recipe
type Recipe struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Cookbook   *Cookbook
}

func (r *Recipe) Type() ObjectType { return RECIPE_OBJ }
func (r *Recipe) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range r.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("recipe")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(r.Body.String())
	out.WriteString("\n}")

	return out.String()
}

// Built In Recipe
type BuiltInRecipe func(args ...Object) Object

type BuiltIn struct {
	Fn BuiltInRecipe
}

func (b *BuiltIn) Type() ObjectType { return BUILT_IN_OBJ }
func (b *BuiltIn) Inspect() string  { return "built-in function" }

// Array
type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// Hash
type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}

	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
