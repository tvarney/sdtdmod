package node

import (
	"math"
	"strconv"
	"strings"

	"github.com/beevik/etree"
	"github.com/tvarney/sdtdmod/pkg/node/key"
)

type Action interface {
	Apply(*etree.Element) bool
	Serialize() map[string]interface{}
}

// Number represents a numeric operation.
//
// This combines a multiply, add, and clamp operation all in one. It also
// contains a precision value used for formatting values back to strings.
type Number struct {
	Attribute string
	Mult      float64
	Add       float64
	Min       float64
	Max       float64
	Precision int
	If        Match
}

// NewNumber returns a Number Action which does nothing to the given attribute.
func NewNumber(attr string) *Number {
	return &Number{
		Attribute: attr,
		Mult:      1.0,
		Add:       0.0,
		Min:       -math.MaxFloat64,
		Max:       math.MaxFloat64,
		Precision: -1,
		If:        nil,
	}
}

// Apply updates the given element.
//
// This will take the attribute in the element, multiply it by `Mult`, add
// `Add`, clamp between `Min` and `Max` inclusive, then truncate to `Precision`
// decimal places.
func (n *Number) Apply(element *etree.Element) bool {
	if n.If != nil && !n.If.Check(element) {
		return false
	}

	attr := element.SelectAttr(n.Attribute)
	if attr == nil {
		return false
	}

	if !strings.Contains(attr.Value, ",") {
		v, err := n.update(attr.Value)
		if err != nil {
			return false
		}
		old := attr.Value
		attr.Value = v
		return old != attr.Value
	}

	builder := strings.Builder{}
	parts := strings.Split(attr.Value, ",")
	v, err := n.update(parts[0])
	if err != nil {
		return false
	}
	builder.WriteString(v)
	for _, p := range parts[1:] {
		v, err = n.update(p)
		if err != nil {
			return false
		}
		builder.WriteRune(',')
		builder.WriteString(v)
	}
	old := attr.Value
	attr.Value = builder.String()
	return (old != attr.Value)
}

func (n *Number) update(value string) (string, error) {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return "", err
	}
	v = math.Max(math.Min((n.Mult*v)+n.Add, n.Max), n.Min)
	s := strconv.FormatFloat(v, 'f', n.Precision, 64)
	if s[len(s)-1] == '0' {
		// Find non-zero end index
		for lastIdx := len(s) - 1; lastIdx > 0; lastIdx-- {
			if s[lastIdx] == '.' {
				// We found the decimal, so trim it
				s = s[:lastIdx]
				if s == "" {
					s = "0"
				}
				break
			}
			if s[lastIdx] != '0' {
				s = s[:lastIdx+1]
				break
			}
		}
	}
	return s, nil
}

// Serialize returns JSON compatible map of this Action.
func (n *Number) Serialize() map[string]interface{} {
	m := map[string]interface{}{
		key.Type: key.ActionNumber,
		"name":   n.Attribute,
	}
	if n.Mult != 1.0 {
		m["mult"] = n.Mult
	}
	if n.Add != 0.0 {
		m["add"] = n.Add
	}
	if n.Min != -math.MaxFloat64 {
		m["min"] = n.Min
	}
	if n.Max != math.MaxFloat64 {
		m["max"] = n.Max
	}
	if n.Precision > 0 {
		m["precision"] = n.Precision
	}
	if n.If != nil {
		m["if"] = n.If.Serialize()
	}

	return m
}

// RemoveAttr is an Action which removes an attribute from the element.
type RemoveAttr struct {
	Attribute string
	If        Match
}

func (r *RemoveAttr) Apply(element *etree.Element) bool {
	if r.If != nil && !r.If.Check(element) {
		return false
	}
	return element.RemoveAttr(r.Attribute) != nil
}

func (r *RemoveAttr) Serialize() map[string]interface{} {
	m := map[string]interface{}{
		key.Type: key.ActionRemoveAttr,
		"name":   r.Attribute,
	}
	if r.If != nil {
		m["if"] = r.If.Serialize()
	}
	return m
}

// InsertAttr is an Action which inserts an attribute into the element.
type InsertAttr struct {
	Attribute string
	Value     string
	If        Match
}

func (i *InsertAttr) Apply(element *etree.Element) bool {
	if i.If != nil && !i.If.Check(element) {
		return false
	}
	element.CreateAttr(i.Attribute, i.Value)
	return true
}

func (i *InsertAttr) Serialize() map[string]interface{} {
	m := map[string]interface{}{
		key.Type: key.ActionInsertAttr,
		"name":   i.Attribute,
		"value":  i.Value,
	}
	if i.If != nil {
		m["if"] = i.If.Serialize()
	}

	return m
}
