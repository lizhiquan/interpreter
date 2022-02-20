package parser

import (
	"fmt"
	"testing"

	"interpreter/ast"
	"interpreter/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("expected len of program.Statements to be 3, got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("expected s.TokenLiteral() to be 'let', got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("expected s to be *ast.LetStatement, got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("expected letStmt.Name.Value to be '%s', got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("expected letStmt.Name.TokenLiteral() to be '%s'. got=%s", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return add(15);
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("expected len of program.Statements to be 3, got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("expected stmt to be a *ast.ReturnStatement, got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("expected returnStmt.TokenLiteral() to be 'return', got=%q", returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 1 {
		t.Fatalf("expected len of program.Statements to be 1, got=%d", len(program.Statements))
	}

	stmt := program.Statements[0]
	exprStmt, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected stmt to be a *ast.ExpressionStatement, got=%T", stmt)
	}

	ident, ok := exprStmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expected exprStmt.Expression to be a *ast.Identifier, got=%T", exprStmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("expected ident.Value to be 'foobar', got=%s", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("expected ident.TokenLiteral() to be 'foobar', got=%s", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 1 {
		t.Fatalf("expected len of program.Statements to be 1, got=%d", len(program.Statements))
	}

	stmt := program.Statements[0]
	exprStmt, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected stmt to be a *ast.ExpressionStatement, got=%T", stmt)
	}

	literal, ok := exprStmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected exprStmt.Expression to be a *ast.IntegerLiteral, got=%T", exprStmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("expected literal.Value to be 5, got=%d", literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("expected ident.TokenLiteral() to be '5', got=%s", literal.TokenLiteral())
	}
}

func TestPrefixExpressions(t *testing.T) {
	tests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}
		if len(program.Statements) != 1 {
			t.Fatalf("expected len of program.Statements to be 1, got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		exprStmt, ok := stmt.(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("expected stmt to be a *ast.ExpressionStatement, got=%T", stmt)
		}

		expr, ok := exprStmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("expected exprStmt.Expression to be a *ast.PrefixExpression, got=%T",
				exprStmt.Expression)
		}
		if expr.Operator != tt.operator {
			t.Errorf("expected expr.Operator to be '%s', got='%s'", tt.operator, expr.Operator)
		}
		if !testIntegerLiteral(t, expr.Right, tt.integerValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, expr ast.Expression, value int64) bool {
	il, ok := expr.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("expected expr to be an *ast.IntegerLiteral, got=%T", expr)
		return false
	}

	if il.Value != value {
		t.Errorf("expected il.Value to be %d, got=%d", value, il.Value)
		return false
	}

	if il.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("expected il.TokenLiteral() to be '%d', got=%s", value, il.TokenLiteral())
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
