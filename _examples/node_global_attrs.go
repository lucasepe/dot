package main

import (
	"fmt"

	"github.com/lucasepe/dot"
)

// go run node_global_attrs.go | dot -Tpng  > node_global_attrs.png

func main() {
	g := dot.NewGraph(dot.Directed)
	g.NodeBaseAttrs().Attr("shape", "plaintext").Attr("color", "blue")
	// Override shape for node `A`
	n1 := g.Node(dot.WithLabel("A"))
	n1.Attr("shape", "box")
	n2 := g.Node(dot.WithLabel("B"))
	n3 := g.Node(dot.WithLabel("C"))

	g.Edge(n1, n2)
	g.Edge(n2, n3).Attr("color", "red")

	fmt.Println(g.String())
}
