package node

import (
	"regexp"
	"strings"

	"github.com/beevik/etree"
	"github.com/tvarney/sdtdmod/pkg/node/key"
)

// Match defines an interface for matching elements.
type Match interface {
	Check(*etree.Element) bool
	Serialize() map[string]interface{}
}

// TagMatch is a match which checks the tag against a set of constraints.
type TagMatch struct {
	Value  string
	Regex  *regexp.Regexp
	Suffix string
	Prefix string
}

func (t *TagMatch) Check(element *etree.Element) bool {
	if t.Value != "" && element.Tag != t.Value {
		return false
	}
	if t.Regex != nil && !t.Regex.MatchString(element.Tag) {
		return false
	}
	if t.Suffix != "" && !strings.HasSuffix(element.Tag, t.Suffix) {
		return false
	}
	if t.Prefix != "" && !strings.HasPrefix(element.Tag, t.Prefix) {
		return false
	}
	return true
}

func (t *TagMatch) Serialize() map[string]interface{} {
	m := map[string]interface{}{
		key.Type: key.MatchTag,
	}
	if t.Value != "" {
		m["value"] = t.Value
	}
	if t.Regex != nil {
		m["regex"] = t.Regex.String()
	}
	if t.Prefix != "" {
		m["prefix"] = t.Prefix
	}
	if t.Suffix != "" {
		m["suffix"] = t.Suffix
	}
	return m
}

// AttrMatch is a match which checks an attribute against a set of constraints.
type AttrMatch struct {
	Attribute string
	Value     string
	Prefix    string
	Suffix    string
	Regex     *regexp.Regexp
}

func (a *AttrMatch) Check(element *etree.Element) bool {
	attr := element.SelectAttr(a.Attribute)
	if attr == nil {
		return false
	}
	if a.Value != "" && attr.Value != a.Value {
		return false
	}
	if a.Regex != nil && !a.Regex.MatchString(attr.Value) {
		return false
	}
	if a.Suffix != "" && !strings.HasSuffix(attr.Value, a.Suffix) {
		return false
	}
	if a.Prefix != "" && !strings.HasPrefix(attr.Value, a.Prefix) {
		return false
	}
	return true
}

func (a *AttrMatch) Serialize() map[string]interface{} {
	m := map[string]interface{}{
		"type": key.MatchAttr,
		"name": a.Attribute,
	}
	if a.Value != "" {
		m["value"] = a.Value
	}
	if a.Regex != nil {
		m["regex"] = a.Regex.String()
	}
	if a.Prefix != "" {
		m["prefix"] = a.Prefix
	}
	if a.Suffix != "" {
		m["suffix"] = a.Suffix
	}
	return m
}

// AnyOf is a Match which checks that an element matches at least one of the
// sub-matches.
type AnyOf []Match

func (a AnyOf) Check(element *etree.Element) bool {
	for _, m := range a {
		if m.Check(element) {
			return true
		}
	}
	return false
}

func (a AnyOf) Serialize() map[string]interface{} {
	matches := make([]map[string]interface{}, 0, len(a))
	for _, m := range a {
		matches = append(matches, m.Serialize())
	}

	return map[string]interface{}{
		"type":    key.MatchAnyOf,
		"matches": matches,
	}
}

// OneOf is a Match which checks that an element matches exactly one of the
// sub-matches.
type OneOf []Match

func (o OneOf) Check(element *etree.Element) bool {
	matches := 0
	for _, m := range o {
		if m.Check(element) {
			matches++
			if matches > 1 {
				return false
			}
		}
	}
	return matches == 1
}

func (o OneOf) Serialize() map[string]interface{} {
	matches := make([]map[string]interface{}, 0, len(o))
	for _, m := range o {
		matches = append(matches, m.Serialize())
	}

	return map[string]interface{}{
		"type":    key.MatchOneOf,
		"matches": matches,
	}
}

// AllOf is a Match which checks that the element matches all of the
// sub-matches.
type AllOf []Match

func (a AllOf) Check(element *etree.Element) bool {
	for _, m := range a {
		if !m.Check(element) {
			return false
		}
	}
	return true
}

func (a AllOf) Serialize() map[string]interface{} {
	matches := make([]map[string]interface{}, 0, len(a))
	for _, m := range a {
		matches = append(matches, m.Serialize())
	}

	return map[string]interface{}{
		"type":    key.MatchAllOf,
		"matches": matches,
	}
}

// Not is a Match which checks that the element does not match the child match.
type Not struct {
	Child Match
}

func (n Not) Check(element *etree.Element) bool {
	return !n.Child.Check(element)
}

func (n Not) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"type":  key.MatchNot,
		"match": n.Child.Serialize(),
	}
}
