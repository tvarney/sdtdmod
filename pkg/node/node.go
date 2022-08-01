package node

import "github.com/beevik/etree"

// Node is a configuration node which may be applied to a xml etree.
type Node struct {
	Match    Match
	Actions  []Action
	Children []*Node
}

// Apply takes an element of an xml etree, checks for matches, and applies if
// matched.
func (n *Node) Apply(element *etree.Element) bool {
	// If we don't match, don't do anything
	if !n.Match.Check(element) {
		return false
	}

	updated := false
	// Iterate over children
	for _, child := range element.ChildElements() {
		for _, node := range n.Children {
			updated = updated || node.Apply(child)
		}
	}

	// Apply any actions
	for _, action := range n.Actions {
		updated = updated || action.Apply(element)
	}
	return updated
}

func (n *Node) Serialize() map[string]interface{} {
	m := map[string]interface{}{}
	if n.Match != nil {
		m["match"] = n.Match.Serialize()
	}
	if len(n.Children) > 0 {
		children := make([]map[string]interface{}, 0, len(n.Children))
		for _, child := range n.Children {
			children = append(children, child.Serialize())
		}
		m["children"] = children
	}
	return m
}
