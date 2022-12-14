package impl

import (
	"fmt"

	"github.com/tvarney/maputil"
	"github.com/tvarney/maputil/errctx"
	"github.com/tvarney/maputil/mpath"
	"github.com/tvarney/maputil/unpack"
	"github.com/tvarney/sdtdmod/pkg/node"
	"github.com/tvarney/sdtdmod/pkg/node/key"
)

// UnpackNodeList takes a JSON list and converts it to a list of Node.
func UnpackNodeList(ctx *errctx.Context, v []interface{}) []*node.Node {
	nodes := make([]*node.Node, 0, len(v))
	for idx, obj := range v {
		ctx.Path.Add(mpath.Index(idx))
		m, ok := obj.(map[string]interface{})
		if !ok {
			ctx.Error(maputil.InvalidTypeError{Expected: []string{maputil.TypeObject}, Actual: maputil.TypeName(obj)})
			ctx.Path.Pop()
			continue
		}

		if n := UnpackNode(ctx, m); n != nil {
			nodes = append(nodes, n)
		}
		ctx.Path.Pop()
	}
	if len(nodes) == 0 {
		return nil
	}
	return nodes
}

// UnpackNode takes a JSON object and unpacks it to a Node.
func UnpackNode(ctx *errctx.Context, v map[string]interface{}) *node.Node {
	rawMatch := unpack.OptionalObject(ctx, v, key.Match, nil)
	rawActions := unpack.OptionalArray(ctx, v, key.Actions, nil)
	rawChildren := unpack.OptionalArray(ctx, v, key.Children, nil)

	n := &node.Node{}

	if rawChildren != nil {
		ctx.Path.Add(mpath.Key(key.Children))
		n.Children = UnpackNodeList(ctx, rawChildren)
		ctx.Path.Pop()
	}
	if rawMatch != nil {
		ctx.Path.Add(mpath.Key(key.Match))
		n.Match = UnpackMatch(ctx, rawMatch)
		ctx.Path.Pop()
	}
	if rawActions != nil {
		ctx.Path.Add(mpath.Key(key.Actions))
		n.Actions = UnpackActionList(ctx, rawActions)
		ctx.Path.Pop()
	}

	if n.Children == nil && n.Actions == nil && n.Match == nil {
		ctx.Error(fmt.Errorf("node requires at least one of %q, %q, or %q", key.Match, key.Actions, key.Children))
		return nil
	}

	return n
}
