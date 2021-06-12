package main

import (
	"fmt"

	"github.com/lucasepe/dot"
)

// go run cluster.go | dot -Tpng  > cluster.png

func main() {
	di := dot.NewGraph(dot.Directed)
	di.Attr("rankdir", "LR")
	outside := di.Node(dot.WithLabel("Outside"))

	// A
	clusterA := di.NewSubgraph()
	clusterA.Attr("label", "Cluster A")

	insideOne := clusterA.Node(dot.WithLabel("one"))
	insideTwo := clusterA.Node(dot.WithLabel("two"))

	// B
	clusterB := di.NewSubgraph()
	clusterB.Attr("label", "Cluster B")

	insideThree := clusterB.Node(dot.WithLabel("three"))
	insideFour := clusterB.Node(dot.WithLabel("four"))

	di.Edge(outside, insideFour)
	di.Edge(insideFour, insideOne)
	di.Edge(insideOne, insideTwo)
	di.Edge(insideTwo, insideThree)
	di.Edge(insideThree, outside)

	fmt.Println(di)
}
