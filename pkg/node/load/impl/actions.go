package impl

import (
	"math"

	"github.com/tvarney/maputil"
	"github.com/tvarney/maputil/errctx"
	"github.com/tvarney/maputil/mpath"
	"github.com/tvarney/maputil/unpack"
	"github.com/tvarney/sdtdmod/pkg/node"
	"github.com/tvarney/sdtdmod/pkg/node/key"
)

// UnpackActionList takes a list of interfaces and unpacks it to an Action
// list.
func UnpackActionList(ctx *errctx.Context, raw []interface{}) []node.Action {
	if len(raw) == 0 {
		return nil
	}

	actions := make([]node.Action, 0, len(raw))
	for idx, v := range raw {
		ctx.Path.Add(mpath.Index(idx))
		obj, err := maputil.AsObject(v)
		if err != nil {
			ctx.Error(err)
			ctx.Path.Pop()
			continue
		}

		action := UnpackAction(ctx, obj)
		if action != nil {
			actions = append(actions, action)
		}
	}

	if len(actions) == 0 {
		return nil
	}
	return actions
}

// UnpackAction takes a JSON object and unpacks it to an Action.
func UnpackAction(ctx *errctx.Context, action map[string]interface{}) node.Action {
	atype := unpack.RequireStringEnum(ctx, action, "type", key.ActionTypes)
	switch atype {
	case key.ActionNumber:
		return UnpackActionNumber(ctx, action)
	case key.ActionInsertAttr:
		return UnpackActionInsertAttr(ctx, action)
	case key.ActionInsertElement:
		return UnpackActionInsertElement(ctx, action)
	case key.ActionRemoveAttr:
		return UnpackActionRemoveAttr(ctx, action)
	case key.ActionRemoveElement:
		return UnpackActionRemoveElement(ctx, action)
	}
	return nil
}

// UnpackActionNumber takes a JSON object and unpacks it to a Number Action.
func UnpackActionNumber(ctx *errctx.Context, action map[string]interface{}) node.Action {
	attr := unpack.RequireString(ctx, action, "attr")
	if attr == "" {
		return nil
	}

	return &node.Number{
		Attribute: attr,
		Mult:      unpack.OptionalNumber(ctx, action, "mult", 1.0),
		Add:       unpack.OptionalNumber(ctx, action, "add", 0.0),
		Min:       unpack.OptionalNumber(ctx, action, "min", -math.MaxFloat64),
		Max:       unpack.OptionalNumber(ctx, action, "max", math.MaxFloat64),
		Precision: int(unpack.OptionalInteger(ctx, action, "precision", -1)),
		If:        UnpackCondition(ctx, action),
	}
}

// UnpackActionInsertAttr takes a JSON object and unpacks it to an InsertAttr
// Action.
func UnpackActionInsertAttr(ctx *errctx.Context, action map[string]interface{}) node.Action {
	errs := ctx.ErrorCount()
	attribute := unpack.RequireString(ctx, action, "name")
	value := unpack.RequireString(ctx, action, "value")
	if ctx.ErrorCount() != errs {
		return nil
	}

	return &node.InsertAttr{
		Attribute: attribute,
		Value:     value,
		If:        UnpackCondition(ctx, action),
	}
}

// UnpackActionInsertElement takes a JSON object and unpacks it to an
// InsertElement Action.
func UnpackActionInsertElement(ctx *errctx.Context, action map[string]interface{}) node.Action {
	return nil
}

// UnpackActionRemoveAttr takes a JSON object and unpacks it to a RemoveAttr
// Action.
func UnpackActionRemoveAttr(ctx *errctx.Context, action map[string]interface{}) node.Action {
	errs := ctx.ErrorCount()
	attribute := unpack.RequireString(ctx, action, "name")
	if ctx.ErrorCount() != errs {
		return nil
	}

	return &node.RemoveAttr{
		Attribute: attribute,
		If:        UnpackCondition(ctx, action),
	}
}

// UnpackActionRemoveElement takes a JSON object and unpacks it to a
// RemoveElement Action.
func UnpackActionRemoveElement(ctx *errctx.Context, action map[string]interface{}) node.Action {
	return nil
}
