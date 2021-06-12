package main

import (
	"fmt"

	"github.com/lucasepe/dot"
)

// go run edges.go | dot -Tpng  > edges.png

func main() {
	g := dot.NewGraph(dot.Directed)
	n1 := g.Node(dot.WithLabel("coding"))
	n2 := g.Node(dot.WithLabel("testing a little"))

	g.Edge(n1, n2)
	g.Edge(n2, n1, "back").Attr("color", "red")

	fmt.Println(g.String())
}
