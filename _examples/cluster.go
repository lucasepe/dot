package main

import (
	"fmt"

	"github.com/lucasepe/dot"
)

// go run cluster.go | dot -Tpng  > cluster.png

func main() {
	di := dot.NewGraph(dot.Directed)
	outside := di.Node("Outside")

	// A
	clusterA := di.Subgraph("Cluster A", dot.ClusterOption{})
	insideOne := clusterA.Node("one")
	insideTwo := clusterA.Node("two")

	// B
	clusterB := di.Subgraph("Cluster B", dot.ClusterOption{})
	insideThree := clusterB.Node("three")
	insideFour := clusterB.Node("four")

	di.Edge(outside, insideFour)
	di.Edge(insideFour, insideOne)
	di.Edge(insideOne, insideTwo)
	di.Edge(insideTwo, insideThree)
	di.Edge(insideThree, outside)

	fmt.Println(di)
}
