package key

const (
	Actions  = "actions"
	Add      = "add"
	Attr     = "attr"
	Children = "children"
	Cond     = "if"
	Match    = "match"
	Matches  = "matches"
	Max      = "max"
	Min      = "min"
	Mult     = "mult"
	Name     = "name"
	Prec     = "precision"
	Prefix   = "prefix"
	Regex    = "regex"
	Suffix   = "suffix"
	Type     = "type"
	Value    = "value"

	MatchTag   = "tag"
	MatchAttr  = "attr"
	MatchAllOf = "all-of"
	MatchAnyOf = "any-of"
	MatchOneOf = "one-of"
	MatchNot   = "not"

	ActionNumber        = "update-number"
	ActionInsertAttr    = "insert-attr"
	ActionInsertElement = "insert-element"
	ActionRemoveAttr    = "remove-attr"
	ActionRemoveElement = "remove-element"
)

var (
	MatchTypes = []string{
		MatchTag, MatchAttr, MatchAllOf, MatchAnyOf, MatchOneOf, MatchNot,
	}
	ActionTypes = []string{
		ActionNumber, ActionInsertAttr, ActionInsertElement, ActionRemoveAttr, ActionRemoveElement,
	}
)
