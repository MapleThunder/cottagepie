package object

import (
	"bytes"
	"cottagepie/ast"
	"fmt"
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

// String
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

// Boolean
type Boolean struct {
	Value bool
}

func (i *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (i *Boolean) Inspect() string  { return fmt.Sprintf("%t", i.Value) }

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
