package main

import (
	p "hook.com/hook/parser"
)

func main() {
	// fs := token.NewFileSet()

	// file, _ := os.ReadFile("parser/ast.go")

	// f, err := parser.ParseFile(fs, "", file, parser.AllErrors)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// p.AddImport(f, "\"hook.com/hook/parser\"", "parser/ast.go")

	// printer.Fprint(os.Stdout, fs, f)

	// // file,_:=os.NewFile()
	// out, _ := os.Create("testFile.txt")

	// printer.Fprint(out, fs, f)
	p.SourceFileAdd("testCases/main.go", "testFile.txt")

}
