package main

import (
	"fmt"

	"github.com/life4/gweb/web"
)

func main() {
	window := web.GetWindow()
	doc := window.Document()
	doc.SetTitle("golangci-lint online")
	issues, err := RunLinters()
	if err != nil {
		panic(err)
	}
	fmt.Println(issues)
	select {}
}
