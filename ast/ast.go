package ast

import "cottagepie/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Program
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// Identifier
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) statementNode()       {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// Bake Statement
type BakeStatement struct {
	Token token.Token // the token.BAKE token
	Name  *Identifier
	Value Expression
}

func (ls *BakeStatement) statementNode()       {}
func (ls *BakeStatement) TokenLiteral() string { return ls.Token.Literal }

// Serves Statement
type ServesStatement struct {
	Token       token.Token // the 'serves' token
	ReturnValue Expression
}

func (rs *ServesStatement) statementNode()       {}
func (rs *ServesStatement) TokenLiteral() string { return rs.Token.Literal }
