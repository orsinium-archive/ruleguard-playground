package main

import (
	"github.com/life4/gweb/web"
)

const Example = `
package main

func main() {
	fmt.Println("hello world")
}
`

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

	rg := NewRuleGuard(editor, &doc)
	rg.Register()
	select {}
}
