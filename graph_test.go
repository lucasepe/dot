package dot

import (
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
	n1 := di.Node(WithLabel("A"))
	n2 := di.Node(WithLabel("B"))
	di.Edge(n1, n2)
	if got, want := flatten(di.String()), `digraph  {n2[label="B"];n1[label="A"];n1->n2;}`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestSubgraph(t *testing.T) {
	di := NewGraph(Directed)
	sub := di.NewSubgraph()
	sub.Attr("label", "test-id").Attr("style", "filled")
	if got, want := flatten(di.String()), `digraph  {subgraph cluster_1 {label="test-id";style="filled";}}`; got != want {
		t.Errorf("got\n[%v] want\n[%v]", got, want)
	}

	sub.Label("hello-sub")
	if got, want := flatten(di.String()), `digraph  {subgraph cluster_1 {label="hello-sub";style="filled";}}`; got != want {
		t.Errorf("got\n[%v] want\n[%v]", got, want)
	}

	found, _ := di.FindSubgraphByLabel("hello-sub")
	if got, want := found, sub; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}

	subsub := sub.NewSubgraph()
	subsub.Attr("label", "hello-sub-sub")

	found, _ = subsub.FindSubgraphByLabel("hello-sub")
	if got, want := found, sub; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}

}

func TestSubgraphClusterOption(t *testing.T) {
	di := NewGraph(Directed)
	sub := di.NewSubgraph()
	sub.Attr("label", "test")
	if got, want := sub.id, "cluster_1"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestEdgeLabel(t *testing.T) {
	di := NewGraph(Directed)
	n1 := di.Node(WithLabel("e1"))
	n2 := di.Node(WithLabel("e2"))
	di.Edge(n1, n2, WithLabel("what"))
	if got, want := flatten(di.String()), `digraph  {n2[label="e2"];n1[label="e1"];n1->n2[label="what"];}`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestSameRank(t *testing.T) {
	di := NewGraph(Directed)
	foo1 := di.Node(WithLabel("foo1"))
	foo2 := di.Node(WithLabel("foo2"))
	bar := di.Node(WithLabel("bar"))

	di.Edge(foo1, foo2)
	di.Edge(foo1, bar)
	di.AddToSameRank("top-row", *foo1, *foo2)

	if got, want := flatten(di.String()), `digraph  {n3[label="bar"];n2[label="foo2"];n1[label="foo1"];n1->n2;n1->n3;{rank=same; n1;n2;};}`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestDeleteLabel(t *testing.T) {
	g := NewGraph()
	n := g.NodeWithID("my-id")
	n.AttributesMap.Delete("label")
	if got, want := flatten(g.String()), `digraph  {n1;}`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestFindNodeByID(t *testing.T) {
	di := NewGraph(Directed)

	if got := di.FindNodeByID("F"); got != nil {
		t.Errorf("got [%v] want [%v]", got, nil)
	}
}

func TestFindNodeByLabel(t *testing.T) {
	di := NewGraph(Directed)
	di.Node(WithLabel("A"))
	di.Node(WithLabel("B"))

	if got := di.FindNodeByLabel("A"); got == nil {
		t.Error("unexpected <nil> value")
	}

	if got := di.FindNodeByLabel("C"); got != nil {
		t.Errorf("got [%v] want [%v]", got, nil)
	}
}

func TestFindNodeByLabelInSubGraphs(t *testing.T) {
	di := NewGraph(Directed)
	di.Node(WithLabel("A"))
	di.Node(WithLabel("B"))
	sub := di.NewSubgraph()
	sub.Attr("label", "new subgraph")
	sub.Node(WithLabel("C"))

	if got := sub.FindNodeByLabel("C"); got == nil {
		t.Errorf("got [%v] want [%v]", nil, "C")
	}
}

func TestLabelWithEscaping(t *testing.T) {
	di := NewGraph(Directed)
	n := di.NodeWithID("without linefeed")
	n.Attr("label", Literal(`"with \l linefeed"`))
	if got, want := flatten(di.String()), `digraph  {n1[label="with \l linefeed"];}`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

// remove tabs and newlines and spaces
func flatten(s string) string {
	return strings.Replace((strings.Replace(s, "\n", "", -1)), "\t", "", -1)
}
