package dot

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

// Node represents a dot Node.
type Node struct {
	AttributesMap
	graph *Graph
	id    string
	seq   int
}

// ID returns the Node identifier
func (n *Node) ID() string {
	return n.id
}

// Seq returns the Node sequential number
func (n *Node) Seq() int {
	return n.seq
}

// Attrs returns the node attributes
func (n *Node) Attrs() *AttributesMap {
	return &n.AttributesMap
}

// Edge represents a graph edge between two Nodes.
type Edge struct {
	AttributesMap
	graph    *Graph
	from, to *Node
}

// Attrs returns the node attributes
func (e *Edge) Attrs() *AttributesMap {
	return &e.AttributesMap
}

// GraphOption is a Graph configuration option.
type GraphOption interface {
	Apply(*Graph)
}

// ClusterOption mark a graph as cluster
type ClusterOption struct{}

// Apply enforces the Graph as cluster
func (o ClusterOption) Apply(g *Graph) {
	g.beCluster()
}

var (
	// Strict defines a `strict` Graph type
	Strict = GraphTypeOption{"strict"}
	// Undirected defines a `graph` Graph type
	Undirected = GraphTypeOption{"graph"}
	// Directed defines a `digraph` Graph type
	Directed = GraphTypeOption{"digraph"}
	// Sub defines a `subgraph` Graph type
	Sub = GraphTypeOption{"subgraph"}
)

// GraphTypeOption sets the graph type
type GraphTypeOption struct {
	Name string
}

// Apply enforces the Graph type
func (o GraphTypeOption) Apply(g *Graph) {
	g.graphType = o.Name
}

// Graph represents a dot graph with nodes and edges.
type Graph struct {
	AttributesMap
	id        string
	isCluster bool
	graphType string
	seq       int
	nodes     map[string]Node
	edgesFrom map[string][]Edge
	subgraphs map[string]*Graph
	parent    *Graph
	sameRank  map[string][]Node
	//
	nodeAttrs AttributesMap
}

// NewGraph return a new initialized Graph.
func NewGraph(options ...GraphOption) *Graph {
	graph := &Graph{
		AttributesMap: AttributesMap{attributes: map[string]interface{}{}},
		graphType:     Directed.Name,
		nodes:         map[string]Node{},
		edgesFrom:     map[string][]Edge{},
		subgraphs:     map[string]*Graph{},
		sameRank:      map[string][]Node{},
		nodeAttrs:     AttributesMap{attributes: map[string]interface{}{}},
	}
	for _, each := range options {
		each.Apply(graph)
	}
	return graph
}

// ID sets the identifier of the graph.
func (g *Graph) ID(newID string) *Graph {
	if len(g.id) > 0 {
		panic("cannot overwrite non-empty id ; both the old and the new could be in use and we cannot tell")
	}
	g.id = newID
	return g
}

// Label sets the "label" attribute value.
func (g *Graph) Label(label string) *Graph {
	g.AttributesMap.Attr("label", label)
	return g
}

// NodeBaseAttrs returns the node global attributes.
func (g *Graph) NodeBaseAttrs() *AttributesMap {
	return &g.nodeAttrs
}

// Root returns the top-level graph if this was a subgraph.
func (g *Graph) Root() *Graph {
	if g.parent == nil {
		return g
	}
	return g.parent.Root()
}

// FindSubgraph returns the subgraph of the graph or one from its parents.
func (g *Graph) FindSubgraph(id string) (*Graph, bool) {
	sub, ok := g.subgraphs[id]
	if !ok {
		if g.parent != nil {
			return g.parent.FindSubgraph(id)
		}
	}
	return sub, ok
}

// FindSubgraphByLabel returns the subgraph of the graph or one from its parents.
func (g *Graph) FindSubgraphByLabel(label string) (*Graph, bool) {
	for _, el := range g.subgraphs {
		if l := el.Value("label"); l != nil {
			if l.(string) == label {
				return el, true
			}
		}
	}

	if g.parent != nil {
		return g.parent.FindSubgraphByLabel(label)
	}

	return nil, false
}

// Subgraph creates a new subgraph.
func (g *Graph) NewSubgraph() *Graph {
	id := fmt.Sprintf("cluster_%d", g.nextSeq())

	sub := NewGraph(Sub)
	sub.id = id
	sub.Attr("label", id) // for consistency with Node creation behavior.
	sub.parent = g
	g.subgraphs[id] = sub
	return sub
}

// Node creates a new node with an autogenerated identifier.
// Eventually specify optional attributes using the `withAttrs` functions.
// The node will have a label attribute with the id as its value.
// Use Label() to overwrite this.
func (g *Graph) Node(withAttrs ...func(*AttributesMap)) *Node {
	return g.NodeWithID("", withAttrs...)
}

// NodeWithID creates a new node with the specified identifier.
// Eventually specify optional attributes using the `withAttrs` functions.
// The node will have a label attribute with the id as its value.
// Use Label() to overwrite this.
func (g *Graph) NodeWithID(id string, withAttrs ...func(*AttributesMap)) *Node {
	seq := g.nextSeq() // create a new, use root sequence

	if len(strings.TrimSpace(id)) == 0 {
		id = fmt.Sprintf("n%d", seq)
	}

	n := Node{
		id:            id,
		seq:           seq,
		AttributesMap: AttributesMap{attributes: map[string]interface{}{"label": id}},
		graph:         g,
	}

	// eventually apply custom attributes
	for _, op := range withAttrs {
		op(n.Attrs())
	}

	// store local
	g.nodes[id] = n

	return &n
}

// Edge creates a new edge between two nodes.
// Eventually specify optional attributes using the `withAttrs` functions.
// Nodes can be have multiple edges to the same other node (or itself).
func (g *Graph) Edge(fromNode, toNode *Node, withAttrs ...func(*AttributesMap)) *Edge {
	// assume fromNode owner == toNode owner
	edgeOwner := g
	if fromNode.graph != toNode.graph { // 1 or 2 are subgraphs
		edgeOwner = commonParentOf(fromNode.graph, toNode.graph)
	}
	e := Edge{
		from:          fromNode,
		to:            toNode,
		AttributesMap: AttributesMap{attributes: map[string]interface{}{}},
		graph:         edgeOwner}

	// eventually apply custom attributes
	for _, op := range withAttrs {
		op(e.Attrs())
	}

	edgeOwner.edgesFrom[fromNode.id] = append(edgeOwner.edgesFrom[fromNode.id], e)
	return &e
}

// FindEdges finds all edges in the graph that go from the fromNode to the toNode.
// Otherwise, returns an empty slice.
func (g *Graph) FindEdges(fromNode, toNode Node) (found []Edge) {
	found = make([]Edge, 0)
	edgeOwner := g
	if fromNode.graph != toNode.graph {
		edgeOwner = commonParentOf(fromNode.graph, toNode.graph)
	}
	if edges, ok := edgeOwner.edgesFrom[fromNode.id]; ok {
		for _, e := range edges {
			if e.to.id == toNode.id {
				found = append(found, e)
			}
		}
	}
	return found
}

// AddToSameRank adds the given nodes to the specified rank group, forcing them to be rendered in the same row
func (g *Graph) AddToSameRank(group string, nodes ...Node) {
	g.sameRank[group] = append(g.sameRank[group], nodes...)
}

// String returns the source in dot notation.
func (g *Graph) String() string {
	b := new(bytes.Buffer)
	g.Write(b)
	return b.String()
}

func (g *Graph) Write(w io.Writer) {
	g.IndentedWrite(NewIndentWriter(w))
}

// IndentedWrite write the graph to a writer using simple TAB indentation.
func (g *Graph) IndentedWrite(w *IndentWriter) {
	fmt.Fprintf(w, "%s %s {", g.graphType, g.id)
	w.NewLineIndentWhile(func() {
		// graph attributes
		if len(g.AttributesMap.attributes) > 0 {
			appendSortedMap(g.AttributesMap.attributes, false, w)
			w.NewLine()
		}

		// node global attributes
		if len(g.nodeAttrs.attributes) > 0 {
			w.NewLine()
			fmt.Fprint(w, "node")
			appendSortedMap(g.nodeAttrs.attributes, true, w)
			w.NewLine()
		}

		// subgraphs
		if len(g.subgraphs) > 0 {
			keys := g.sortedSubgraphsKeys()

			for _, key := range keys {
				w.NewLine()
				if each, ok := g.FindSubgraph(key); ok {
					//each := g.subgraphs[key]
					each.IndentedWrite(w)
				}
			}
		}

		// nodes
		if tot := len(g.nodes); tot > 0 {
			w.NewLine()

			nodeKeys := g.sortedNodesKeys()

			for i, key := range nodeKeys {
				each := g.nodes[key]
				fmt.Fprintf(w, "n%d", each.seq)
				appendSortedMap(each.attributes, true, w)
				fmt.Fprintf(w, ";")
				if i < tot-1 {
					w.NewLine()
				}
			}
		}

		// edges
		if tot := len(g.edgesFrom); tot > 0 {
			w.NewLine()
			w.NewLine()

			denoteEdge := "->"
			if g.graphType == "graph" {
				denoteEdge = "--"
			}

			edgeKeys := g.sortedEdgesFromKeys()

			for i, each := range edgeKeys {
				all := g.edgesFrom[each]
				for _, each := range all {
					fmt.Fprintf(w, "n%d%sn%d", each.from.seq, denoteEdge, each.to.seq)
					appendSortedMap(each.attributes, true, w)
					fmt.Fprint(w, ";")
					if i < tot-1 {
						w.NewLine()
					}
				}
			}
		}

		if tot := len(g.sameRank); tot > 0 {
			w.NewLine()

			for _, nodes := range g.sameRank {
				str := ""
				for _, n := range nodes {
					str += fmt.Sprintf("n%d;", n.seq)
				}
				fmt.Fprintf(w, "{rank=same; %s};", str)
				w.NewLine()
			}
		}
	})

	fmt.Fprintf(w, "}")
	w.NewLine()
}

// VisitNodes visits all nodes recursively
func (g *Graph) VisitNodes(callback func(node *Node) (done bool)) {
	for _, node := range g.nodes {
		done := callback(&node)
		if done {
			return
		}
	}

	for _, subGraph := range g.subgraphs {
		subGraph.VisitNodes(callback)
	}
}

// FindNodeByID returns a node by its identifier.
func (g *Graph) FindNodeByID(id string) (found *Node) {
	g.VisitNodes(func(node *Node) (done bool) {
		if node.id == id {
			found = node
			return true
		}
		return false
	})
	return
}

// FindNodeByLabel returns a node by its label.
func (g *Graph) FindNodeByLabel(label string) (found *Node) {
	g.VisitNodes(func(node *Node) (done bool) {
		if l := node.Value("label"); l != nil {
			if l.(string) == label {
				found = node
				return true
			}
		}
		return false
	})
	return
}

func (g *Graph) beCluster() {
	g.id = "cluster_" + g.id
}

func commonParentOf(one *Graph, two *Graph) *Graph {
	// TODO
	return one.Root()
}

func appendSortedMap(m map[string]interface{}, mustBracket bool, b io.Writer) {
	if len(m) == 0 {
		return
	}
	if mustBracket {
		fmt.Fprint(b, "[")
	}
	first := true
	// first collect keys
	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}
	sort.StringSlice(keys).Sort()

	for _, k := range keys {
		if !first {
			if mustBracket {
				fmt.Fprint(b, ",")
			} else {
				fmt.Fprintf(b, ";")
			}
		}
		if html, isHTML := m[k].(HTML); isHTML {
			fmt.Fprintf(b, "%s=<%s>", k, html)
		} else if literal, isLiteral := m[k].(Literal); isLiteral {
			fmt.Fprintf(b, "%s=%s", k, literal)
		} else {
			fmt.Fprintf(b, "%s=%q", k, m[k])
		}
		first = false
	}
	if mustBracket {
		fmt.Fprint(b, "]")
	} else {
		fmt.Fprint(b, ";")
	}
}

func (g *Graph) sortedNodesKeys() (keys []string) {
	for each := range g.nodes {
		keys = append(keys, each)
	}

	sort.Slice(keys, func(i, j int) bool {
		x, _ := strconv.Atoi(keys[i][1:])
		y, _ := strconv.Atoi(keys[j][1:])
		return x > y
	})
	return
}

func (g *Graph) sortedEdgesFromKeys() (keys []string) {
	for each := range g.edgesFrom {
		keys = append(keys, each)
	}

	sort.Slice(keys, func(i, j int) bool {
		x, _ := strconv.Atoi(keys[i][1:])
		y, _ := strconv.Atoi(keys[j][1:])
		return x > y
	})

	return
}

func (g *Graph) sortedSubgraphsKeys() (keys []string) {
	for _, v := range g.subgraphs {
		keys = append(keys, v.id)
	}

	sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	return
}

// nextSeq takes the next sequence number from the root graph
func (g *Graph) nextSeq() int {
	root := g.Root()
	root.seq++
	return root.seq
}
