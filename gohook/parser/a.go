package main

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
)

func main() {
	// src is the input for which we want to print the AST.
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "t.go", nil, 0)
	if err != nil {
		panic(err)
	}

	// Print the AST.
	ast.Print(fset, f)
	isInsert := false
	var mainBody *ast.BlockStmt
	var toInsertIfStmt *ast.IfStmt
	for _, decl := range f.Decls {

		fc, ok := decl.(*ast.FuncDecl)

		if ok && fc.Name.Name == "useless" {
			ifstmt, okIf := fc.Body.List[0].(*ast.IfStmt)
			if okIf {
				print("ok")
				toInsertIfStmt = ifstmt
			}
		}
		if ok && fc.Name.Name == "main" {
			mainBody = fc.Body
		}
	}
	if !isInsert && mainBody != nil && toInsertIfStmt != nil {
		t := make([]ast.Stmt, 0, len(mainBody.List)+1)
		// toInsertIfStmt = &ast.IfStmt{
		// 	If:   toInsertIfStmt.If,
		// 	Init: toInsertIfStmt.Init,
		// 	Cond: toInsertIfStmt.Cond,
		// 	Body: toInsertIfStmt.Body,
		// 	Else: toInsertIfStmt.Else,
		// }
		t = append(t, toInsertIfStmt)
		t = append(t, mainBody.List...)
		mainBody.List = t
	}
	ast.Print(fset, f)
	var cfg printer.Config
	cfg.Fprint(os.Stderr, fset, f)

}
