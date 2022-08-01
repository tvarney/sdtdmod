package key

const (
	Type = "type"
	Cond = "if"

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
