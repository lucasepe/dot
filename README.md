## dot - little helper package in Go for the graphviz dot language

[![Go Report Card](https://goreportcard.com/badge/github.com/lucasepe/dot)](https://goreportcard.com/report/github.com/lucasepe/dot)
[![GoDoc](https://godoc.org/github.com/lucasepe/dot?status.svg)](https://pkg.go.dev/github.com/lucasepe/dot)

This is a modified fork of the [great project](https://github.com/emicklei/dot/) made by [Ernest Micklei](https://github.com/emicklei). 


[DOT language](http://www.graphviz.org/doc/info/lang.html)

	package main
	
	import (
		"fmt"	
		"github.com/lucasepe/dot"
	)
	
	// go run main.go | dot -Tpng  > test.png && open test.png
	
	func main() {
		g := dot.NewGraph(dot.Directed)
		n1 := g.Node("coding")
		n2 := g.Node("testing a little").Box()
	
		g.Edge(n1, n2)
		g.Edge(n2, n1, "back").Attr("color", "red")
	
		fmt.Println(g.String())
	}

Output

	digraph {
		node [label="coding"]; n1;
		node [label="testing a little"];
		n1 -> n2;
		n2 -> n1 [color="red", label="back"];
	}

Subgraphs

	s := g.Subgraph("cluster")
	s.Attr("style","filled")


Initializers

	g := dot.NewGraph(dot.Directed)
	g.NodeInitializer(func(n dot.Node) {
		n.Attr("shape", "rectangle")
		n.Attr("fontname", "arial")
		n.Attr("style", "rounded,filled")
	})

	g.EdgeInitializer(func(e dot.Edge) {
		e.Attr("fontname", "arial")
		e.Attr("fontsize", "9")
		e.Attr("arrowsize", "0.8")
		e.Attr("arrowhead", "open")
	})

HTML and Literal values

	node.Attr("label", Literal(`"left-justified text\l"`))
	graph.Attr("label", HTML("<B>Hi</B>"))

Nodes Global Attributes

    g := dot.NewGraph(dot.Directed)
	g.NodeGlobalAttrs("shape", "plaintext", "color", "blue")
	// Override shape for node `A`
	n1 := g.Node("A").Attr("shape", "box")

## cluster example

![](./_examples/cluster.png)

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

	// edges
	di.Edge(outside, insideFour)
	di.Edge(insideFour, insideOne)
	di.Edge(insideOne, insideTwo)
	di.Edge(insideTwo, insideThree)
	di.Edge(insideThree, outside)

## About dot attributes

https://graphviz.gitlab.io/_pages/doc/info/attrs.html

## display your graph

	go run main.go | dot -Tpng  > test.png && open test.png


