package dot

import (
	"fmt"
	"io"
	"sort"
)

// HTML renders the provided content as graphviz HTML. Use of this
// type is only valid for some attributes, like the 'label' attribute.
type HTML string

// Literal renders the provided value as is, without adding enclosing
// quotes, escaping newlines, quotations marks or any other characters.
// For example:
//     node.Attr("label", Literal(`"left-justified text\l"`))
// allows you to left-justify the label (due to the \l at the end).
// The caller is responsible for enclosing the value in quotes and for
// proper escaping of special characters.
type Literal string

// AttributesMap holds attribute=value pairs.
type AttributesMap struct {
	attributes map[string]interface{}
}

// Attr sets the value for an attribute (unless empty).
func (a *AttributesMap) Attr(label string, value interface{}) *AttributesMap {
	if len(label) == 0 || value == nil {
		return a
	}

	if s, ok := value.(string); ok {
		if len(s) > 0 {
			a.attributes[label] = s
			return a
		}
	}

	a.attributes[label] = value
	return a
}

// Value return the value added for this label.
func (a *AttributesMap) Value(label string) interface{} {
	return a.attributes[label]
}

// Delete removes the attribute value at key, if any
func (a *AttributesMap) Delete(key string) {
	delete(a.attributes, key)
}

func (a *AttributesMap) Write(wri io.Writer, mustBracket bool) {
	if len(a.attributes) == 0 {
		return
	}

	if mustBracket {
		fmt.Fprint(wri, "[")
	}
	first := true
	// first collect keys
	keys := []string{}
	for k := range a.attributes {
		keys = append(keys, k)
	}
	sort.StringSlice(keys).Sort()

	for _, k := range keys {
		if !first {
			if mustBracket {
				fmt.Fprint(wri, ",")
			} else {
				fmt.Fprintf(wri, ";")
			}
		}
		if html, isHTML := a.attributes[k].(HTML); isHTML {
			fmt.Fprintf(wri, "%s=<%s>", k, html)
		} else if literal, isLiteral := a.attributes[k].(Literal); isLiteral {
			fmt.Fprintf(wri, "%s=%s", k, literal)
		} else {
			fmt.Fprintf(wri, "%s=%q", k, a.attributes[k])
		}
		first = false
	}
	if mustBracket {
		fmt.Fprint(wri, "]")
	} else {
		fmt.Fprint(wri, ";")
	}
}
