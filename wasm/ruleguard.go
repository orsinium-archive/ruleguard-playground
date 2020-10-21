package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"

	"github.com/quasilyte/go-ruleguard/analyzer"
	"golang.org/x/tools/go/analysis"
)

func RunRuleGuard(src string, rule string) ([]analysis.Diagnostic, error) {
	const filename = "main.go"
	const mode = parser.AllErrors | parser.ParseComments

	issues := make([]analysis.Diagnostic, 0)
	fset := token.NewFileSet()
	tree, err := parser.ParseFile(fset, filename, src, mode)
	if err != nil {
		return nil, err
	}
	pass := analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{tree},
		Pkg:   types.NewPackage("main", "main"),
		Report: func(d analysis.Diagnostic) {
			issues = append(issues, d)
		},
	}

	an := analyzer.Analyzer
	an.Flags.Set("e", rule)
	_, err = an.Run(&pass)
	if err != nil {
		return nil, err
	}
	return issues, nil
}
