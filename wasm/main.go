package main

import (
	"fmt"
	"syscall/js"

	"github.com/life4/gweb/web"
)

const Example = `
package main

func main() {
	fmt.Println("hello world")
}
`

func lint(editor web.Value, doc *web.Document) {
	out := doc.Element("lint-output")
	rule := doc.Element("lint-rule").Get("value").String()
	src := editor.Call("getValue").String()

	issues, err := RunRuleGuard(src, rule)
	if err != nil {
		out.SetText(err.Error())
		return
	}
	fmt.Println(issues)
	out.SetText(fmt.Sprintln(issues))
}

func register(editor web.Value, doc *web.Document) {
	btn := doc.Element("lint-run")
	btn.Set("disabled", false)

	wrapped := func(this js.Value, args []js.Value) interface{} {
		btn.Set("disabled", true)
		lint(editor, doc)
		register(editor, doc)
		return true
	}
	btn.Call("addEventListener", "click", js.FuncOf(wrapped))
}

func main() {
	window := web.GetWindow()
	doc := window.Document()
	doc.SetTitle("golangci-lint online")

	// init code editor
	input := doc.Element("lint-code")
	input.SetInnerHTML(Example)
	editor := window.Get("CodeMirror").Call("fromTextArea",
		input,
		map[string]interface{}{
			"lineNumbers": true,
		})

	register(editor, &doc)
	select {}
}
