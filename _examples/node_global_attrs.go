package main

import (
	"fmt"

	"github.com/lucasepe/dot"
)

// go run node_global_attrs.go | dot -Tpng  > node_global_attrs.png

func main() {
	g := dot.NewGraph(dot.Directed)
	g.NodeGlobalAttrs("shape", "plaintext", "color", "blue")
	// Override shape for node `A`
	n1 := g.Node("A").Attr("shape", "box")
	n2 := g.Node("B")
	n3 := g.Node("C")

	g.Edge(n1, n2)
	g.Edge(n2, n3).Attr("color", "red")

	fmt.Println(g.String())
}
