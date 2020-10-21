package main

import (
	"fmt"

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

	// run linter
	src := editor.Call("getValue").String()
	rule := "m.Match(`$x`)"
	issues, err := RunRuleGuard(src, rule)
	if err != nil {
		panic(err)
	}
	fmt.Println(issues)
	select {}
}
