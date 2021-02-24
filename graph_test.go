package dot

import (
	"reflect"
	"strings"
	"testing"
)

func TestEmpty(t *testing.T) {
	di := NewGraph(Directed)
	if got, want := flatten(di.String()), `digraph  {}`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestEmptyWithIDAndAttributes(t *testing.T) {
	di := NewGraph(Directed)
	di.ID("test")
	di.Attr("style", "filled")
	di.Attr("color", "lightgrey")
	if got, want := flatten(di.String()), `digraph test {color="lightgrey";style="filled";}`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestEmptyWithHTMLLabel(t *testing.T) {
	di := NewGraph(Directed)
	di.ID("test")
	di.Attr("label", HTML("<B>Hi</B>"))
	if got, want := flatten(di.String()), `digraph test {label=<<B>Hi</B>>;}`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestEmptyWithLiteralValueLabel(t *testing.T) {
	di := NewGraph(Directed)
	di.ID("test")
	di.Attr("label", Literal(`"left-justified text\l"`))
	if got, want := flatten(di.String()), `digraph test {label="left-justified text\l";}`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestTwoConnectedNodes(t *testing.T) {
	di := NewGraph(Directed)
	n1 := di.Node("A")
	n2 := di.Node("B")
	di.Edge(n1, n2)
	if got, want := flatten(di.String()), `digraph  {n1[label="A"];n2[label="B"];n1->n2;}`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestGraph_FindEdges(t *testing.T) {
	di := NewGraph(Directed)
	n1 := di.Node("A")
	n2 := di.Node("B")
	want := []Edge{di.Edge(n1, n2)}
	got := di.FindEdges(n1, n2)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("TestGraph.FindEdges() = %v, want %v", got, want)
	}
}

func TestSubgraph(t *testing.T) {
	di := NewGraph(Directed)
	sub := di.Subgraph("test-id")
	sub.Attr("style", "filled")
	if got, want := flatten(di.String()), `digraph  {subgraph s1 {label="test-id";style="filled";}}`; got != want {
		t.Errorf("got\n[%v] want\n[%v]", got, want)
	}
	sub.Label("new-label")
	if got, want := flatten(di.String()), `digraph  {subgraph s1 {label="new-label";style="filled";}}`; got != want {
		t.Errorf("got\n[%v] want\n[%v]", got, want)
	}
	found, _ := di.FindSubgraph("test-id")
	if got, want := found, sub; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	subsub := sub.Subgraph("sub-test-id")
	found, _ = subsub.FindSubgraph("test-id")
	if got, want := found, sub; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}

}

func TestSubgraphClusterOption(t *testing.T) {
	di := NewGraph(Directed)
	sub := di.Subgraph("test", ClusterOption{})
	if got, want := sub.id, "cluster_s1"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestEdgeLabel(t *testing.T) {
	di := NewGraph(Directed)
	n1 := di.Node("e1")
	n2 := di.Node("e2")
	di.Edge(n1, n2, "what")
	if got, want := flatten(di.String()), `digraph  {n1[label="e1"];n2[label="e2"];n1->n2[label="what"];}`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestSameRank(t *testing.T) {
	di := NewGraph(Directed)
	foo1 := di.Node("foo1")
	foo2 := di.Node("foo2")
	bar := di.Node("bar")

	di.Edge(foo1, foo2)
	di.Edge(foo1, bar)
	di.AddToSameRank("top-row", foo1, foo2)

	if got, want := flatten(di.String()), `digraph  {n1[label="foo1"];n2[label="foo2"];n3[label="bar"];n1->n2;n1->n3;{rank=same; n1;n2;};}`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestDeleteLabel(t *testing.T) {
	g := NewGraph()
	n := g.Node("my-id")
	n.AttributesMap.Delete("label")
	if got, want := flatten(g.String()), `digraph  {n1;}`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestGraph_FindNodeById_emptyGraph(t *testing.T) {
	di := NewGraph(Directed)

	_, found := di.FindNode("F")

	if got, want := found, false; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestGraph_FindNodeByLabel(t *testing.T) {
	di := NewGraph(Directed)
	di.Node("A")
	di.Node("B")

	node, found := di.FindNode("A")

	if got, want := node.label, "A"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}

	if got, want := found, true; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestGraph_FindNodeInSubGraphs(t *testing.T) {
	di := NewGraph(Directed)
	di.Node("A")
	di.Node("B")
	sub := di.Subgraph("new subgraph")
	sub.Node("C")

	node, found := sub.FindNode("C")

	if got, want := node.label, "C"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}

	if got, want := found, true; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestGraph_FindNodesInSubGraphs(t *testing.T) {
	di := NewGraph(Directed)
	di.Node("A")
	di.Node("B")
	sub := di.Subgraph("new subgraph")
	sub.Node("C")

	nodes := di.FindNodes()

	if got, want := len(nodes), 3; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestLabelWithEscaping(t *testing.T) {
	di := NewGraph(Directed)
	n := di.Node("without linefeed")
	n.Attr("label", Literal(`"with \l linefeed"`))
	if got, want := flatten(di.String()), `digraph  {n1[label="with \l linefeed"];}`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestGraphCreateNodeOnce(t *testing.T) {
	di := NewGraph(Undirected)
	n1 := di.Node("A")
	n2 := di.Node("A")
	if got, want := n1, n2; &n1 == &n2 {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

// remove tabs and newlines and spaces
func flatten(s string) string {
	return strings.Replace((strings.Replace(s, "\n", "", -1)), "\t", "", -1)
}
