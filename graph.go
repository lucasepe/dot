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
	label string
}

// ID returns the Node identifier
func (n Node) ID() string {
	return n.id
}

// Seq returns the Node sequential number
func (n Node) Seq() int {
	return n.seq
}

// Attr sets label=value and return the Node
func (n Node) Attr(label string, value interface{}) Node {
	n.AttributesMap.Attr(label, value)
	return n
}

// Edge represents a graph edge between two Nodes.
type Edge struct {
	AttributesMap
	graph    *Graph
	from, to *Node
}

// Attr sets key=value and returns the Egde.
func (e Edge) Attr(key string, value interface{}) Edge {
	e.AttributesMap.Attr(key, value)
	return e
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
	nodes     map[string]*Node
	edgesFrom map[string][]Edge
	subgraphs map[string]*Graph
	parent    *Graph
	sameRank  map[string][]*Node
	//
	nodeGlobalAttrs AttributesMap
}

// NewGraph return a new initialized Graph.
func NewGraph(options ...GraphOption) *Graph {
	graph := &Graph{
		AttributesMap:   AttributesMap{attributes: map[string]interface{}{}},
		graphType:       Directed.Name,
		nodes:           map[string]*Node{},
		edgesFrom:       map[string][]Edge{},
		subgraphs:       map[string]*Graph{},
		sameRank:        map[string][]*Node{},
		nodeGlobalAttrs: AttributesMap{attributes: map[string]interface{}{}},
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

// NodeGlobalAttrs sets the global attributes for all nodes.
func (g *Graph) NodeGlobalAttrs(labelvalues ...interface{}) {
	if len(labelvalues)%2 != 0 {
		panic("missing label or value ; must provide pairs")
	}
	for i := 0; i < len(labelvalues); i += 2 {
		label := labelvalues[i].(string)
		value := labelvalues[i+1]
		g.nodeGlobalAttrs.Attr(label, value)
	}
}

func (g *Graph) beCluster() {
	g.id = "cluster_" + g.id
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

// Subgraph returns the Graph with the given id ; creates one if absent.
// The label attribute is also set to the id ; use Label() to overwrite it.
func (g *Graph) Subgraph(id string, options ...GraphOption) *Graph {
	sub, ok := g.subgraphs[id]
	if ok {
		return sub
	}
	sub = NewGraph(Sub)
	sub.Attr("label", id) // for consistency with Node creation behavior.
	sub.id = fmt.Sprintf("s%d", g.nextSeq())
	for _, each := range options {
		each.Apply(sub)
	}
	sub.parent = g
	g.subgraphs[id] = sub
	return sub
}

// FindNode finds a Node by it's label.
func (g *Graph) FindNode(label string) (*Node, bool) {
	for _, v := range g.nodes {
		if v.label == label {
			return v, true
		}
	}
	if g.parent == nil {
		return &Node{}, false
	}
	return g.parent.FindNode(label)
}

// nextSeq takes the next sequence number from the root graph
func (g *Graph) nextSeq() int {
	root := g.Root()
	root.seq++
	return root.seq
}

// Node returns the node created with this id or creates a new node if absent.
// The node will have a label attribute with the id as its value. Use Label() to overwrite this.
// This method can be used as both a constructor and accessor.
// not thread safe!
//
// Node creates a new node.
// Warning: if a Node with this identifier already exists it will be overwritten.
func (g *Graph) Node(label string) *Node {
	seq := g.nextSeq() // create a new, use root sequence
	id := fmt.Sprintf("n%d", seq)
	n := &Node{
		id:    id,
		seq:   seq,
		label: label,
		AttributesMap: AttributesMap{attributes: map[string]interface{}{
			"label": label}},
		graph: g,
	}

	// store local
	g.nodes[id] = n
	return n
}

// Edge creates a new edge between two nodes.
// Nodes can be have multiple edges to the same other node (or itself).
// If one or more labels are given then the "label" attribute is set to the edge.
func (g *Graph) Edge(fromNode, toNode *Node, labels ...string) Edge {
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
	if len(labels) > 0 {
		e.Attr("label", strings.Join(labels, ","))
	}

	edgeOwner.edgesFrom[fromNode.id] = append(edgeOwner.edgesFrom[fromNode.id], e)
	return e
}

// FindEdges finds all edges in the graph that go from the fromNode to the toNode.
// Otherwise, returns an empty slice.
func (g *Graph) FindEdges(fromNode, toNode *Node) (found []Edge) {
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

func commonParentOf(one *Graph, two *Graph) *Graph {
	// TODO
	return one.Root()
}

// AddToSameRank adds the given nodes to the specified rank group, forcing them to be rendered in the same row
func (g *Graph) AddToSameRank(group string, nodes ...*Node) {
	g.sameRank[group] = append(g.sameRank[group], nodes...)
}

// String returns the source in dot notation.
func (g Graph) String() string {
	b := new(bytes.Buffer)
	g.Write(b)
	return b.String()
}

func (g Graph) Write(w io.Writer) {
	g.IndentedWrite(NewIndentWriter(w))
}

// IndentedWrite write the graph to a writer using simple TAB indentation.
func (g Graph) IndentedWrite(w *IndentWriter) {
	fmt.Fprintf(w, "%s %s {", g.graphType, g.id)
	w.NewLineIndentWhile(func() {
		// subgraphs
		for _, key := range g.sortedSubgraphsKeys() {
			each := g.subgraphs[key]
			each.IndentedWrite(w)
		}
		// graph attributes
		g.AttributesMap.Write(w, false)
		w.NewLine()

		// graph nodes global attributes
		if len(g.nodeGlobalAttrs.attributes) > 0 {
			fmt.Fprintf(w, "node ")
			g.nodeGlobalAttrs.Write(w, true)
			fmt.Fprintf(w, ";")
			w.NewLine()
		}
		w.NewLine()
		// graph nodes
		for _, key := range g.sortedNodesKeys() {
			each := g.nodes[key]
			fmt.Fprintf(w, "n%d", each.seq)
			each.AttributesMap.Write(w, true)
			fmt.Fprintf(w, ";")
			w.NewLine()
		}
		w.NewLine()

		// graph edges
		denoteEdge := "->"
		if g.graphType == "graph" {
			denoteEdge = "--"
		}
		for _, each := range g.sortedEdgesFromKeys() {
			all := g.edgesFrom[each]
			for _, each := range all {
				fmt.Fprintf(w, "n%d%sn%d", each.from.seq, denoteEdge, each.to.seq)
				each.AttributesMap.Write(w, true)
				fmt.Fprint(w, ";")
				w.NewLine()
			}
		}
		for _, nodes := range g.sameRank {
			str := ""
			for _, n := range nodes {
				str += fmt.Sprintf("n%d;", n.seq)
			}
			fmt.Fprintf(w, "{rank=same; %s};", str)
			w.NewLine()
		}
	})
	fmt.Fprintf(w, "}")
	w.NewLine()
}

// VisitNodes visits all nodes recursively
func (g Graph) VisitNodes(callback func(node *Node) (done bool)) {
	for _, node := range g.nodes {
		done := callback(node)
		if done {
			return
		}
	}

	for _, subGraph := range g.subgraphs {
		subGraph.VisitNodes(callback)
	}
}

// FindNodes returns all nodes recursively
func (g Graph) FindNodes() (nodes []*Node) {
	var foundNodes []*Node
	g.VisitNodes(func(node *Node) (done bool) {
		foundNodes = append(foundNodes, node)
		return false
	})
	return foundNodes
}

func (g *Graph) sortedNodesKeys() (keys []string) {
	for each := range g.nodes {
		keys = append(keys, each)
	}

	sort.Slice(keys, func(i, j int) bool {
		x, _ := strconv.Atoi(keys[i][1:])
		y, _ := strconv.Atoi(keys[j][1:])
		return x < y
	})
	//sort.StringSlice(keys).Sort()
	return
}
func (g *Graph) sortedEdgesFromKeys() (keys []string) {
	for each := range g.edgesFrom {
		keys = append(keys, each)
	}
	sort.StringSlice(keys).Sort()
	return
}
func (g *Graph) sortedSubgraphsKeys() (keys []string) {
	for each := range g.subgraphs {
		keys = append(keys, each)
	}
	sort.StringSlice(keys).Sort()
	return
}
