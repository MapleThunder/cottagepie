package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	SERVES_VALUE_OBJ = "SERVES_VALUE"
	ERROR_OBJ        = "ERROR"
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
