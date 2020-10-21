package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"syscall/js"

	"github.com/life4/gweb/web"
	"github.com/quasilyte/go-ruleguard/analyzer"
	"golang.org/x/tools/go/analysis"
)

const filename = "main.go"
const mode = parser.AllErrors | parser.ParseComments

type RuleGuard struct {
	out    web.HTMLElement
	editor web.Value
	doc    *web.Document
}

func NewRuleGuard(editor web.Value, doc *web.Document) RuleGuard {
	return RuleGuard{
		editor: editor,
		doc:    doc,
		out:    doc.Element("lint-output"),
	}
}

func (rg *RuleGuard) Register() {
	btn := rg.doc.Element("lint-run")
	btn.Set("disabled", false)

	wrapped := func(this js.Value, args []js.Value) interface{} {
		btn.Set("disabled", true)
		rg.RunAndPrint()
		rg.Register()
		return true
	}
	btn.Call("addEventListener", "click", js.FuncOf(wrapped))
}

func (rg *RuleGuard) RunAndPrint() {
	rg.out.SetText("Running...")
	rule := rg.doc.Element("lint-rule").Get("value").String()
	src := rg.editor.Call("getValue").String()

	// parse source code
	pass, err := rg.makePass(src)
	if err != nil {
		rg.out.SetText(err.Error())
		return
	}

	// run linter
	issues, err := rg.Run(pass, src, rule)
	if err != nil {
		rg.out.SetText(err.Error())
		return
	}

	// print violations
	rg.table(issues, pass)
}

func (rg *RuleGuard) makePass(src string) (*analysis.Pass, error) {
	fset := token.NewFileSet()
	tree, err := parser.ParseFile(fset, filename, src, mode)
	if err != nil {
		return nil, err
	}
	pass := analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{tree},
		Pkg:   types.NewPackage("main", "main"),
	}
	return &pass, nil
}
func (rg *RuleGuard) Run(pass *analysis.Pass, src string, rule string) ([]analysis.Diagnostic, error) {
	issues := make([]analysis.Diagnostic, 0)
	pass.Report = func(d analysis.Diagnostic) {
		issues = append(issues, d)
	}
	an := analyzer.Analyzer
	err := an.Flags.Set("e", rule)
	if err != nil {
		return nil, err
	}
	_, err = an.Run(pass)
	if err != nil {
		return nil, err
	}
	return issues, nil
}

func (rg *RuleGuard) table(issues []analysis.Diagnostic, pass *analysis.Pass) {
	table := rg.doc.CreateElement("table")
	table.Attribute("class").Set("table table-sm")

	thead := rg.doc.CreateElement("thead")
	table.Node().AppendChild(thead.Node())
	tr := rg.doc.CreateElement("tr")
	thead.Node().AppendChild(tr.Node())

	cols := []string{"pos", "message"}
	for _, name := range cols {
		th := rg.doc.CreateElement("th")
		th.SetText(name)
		tr.Node().AppendChild(th.Node())
	}

	tbody := rg.doc.CreateElement("tbody")
	table.Node().AppendChild(tbody.Node())

	for _, issue := range issues {
		tr := rg.doc.CreateElement("tr")

		td := rg.doc.CreateElement("td")
		pos := pass.Fset.Position(issue.Pos)
		td.SetText(fmt.Sprintf("%d:%d", pos.Line, pos.Column))
		tr.Node().AppendChild(td.Node())

		td = rg.doc.CreateElement("td")
		td.SetText(issue.Message)
		tr.Node().AppendChild(td.Node())

		tbody.Node().AppendChild(tr.Node())
	}

	rg.out.Node().AppendChild(table.Node())
}
