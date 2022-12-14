package impl

import (
	"github.com/tvarney/maputil"
	"github.com/tvarney/maputil/errctx"
	"github.com/tvarney/maputil/mpath"
	"github.com/tvarney/maputil/unpack"
	"github.com/tvarney/sdtdmod/pkg/node"
	"github.com/tvarney/sdtdmod/pkg/node/key"
)

func UnpackMatchArray(ctx *errctx.Context, match map[string]interface{}) []node.Match {
	raw := unpack.RequireArray(ctx, match, key.Matches)
	if len(raw) == 0 {
		return nil
	}
	ctx.Path.Add(mpath.Key(key.Matches))

	matches := make([]node.Match, 0, len(raw))
	for idx, v := range raw {
		obj, err := maputil.AsObject(v)
		if err != nil {
			ctx.ErrorWithIndex(err, idx)
			continue
		}

		ctx.Path.Add(mpath.Index(idx))
		m := UnpackMatch(ctx, obj)
		ctx.Path.Pop()
		if m != nil {
			matches = append(matches, m)
		}
	}
	ctx.Path.Pop()
	if len(matches) == 0 {
		return nil
	}
	return matches
}
func UnpackMatch(ctx *errctx.Context, match map[string]interface{}) node.Match {
	mtype := unpack.RequireStringEnum(ctx, match, key.Type, key.MatchTypes)
	switch mtype {
	case key.MatchTag:
		return UnpackMatchTag(ctx, match)
	case key.MatchAttr:
		return UnpackMatchAttr(ctx, match)
	case key.MatchAllOf:
		return UnpackMatchAllOf(ctx, match)
	case key.MatchAnyOf:
		return UnpackMatchAnyOf(ctx, match)
	case key.MatchOneOf:
		return UnpackMatchOneOf(ctx, match)
	case key.MatchNot:
		return UnpackMatchNot(ctx, match)
	}
	return nil
}

func UnpackMatchTag(ctx *errctx.Context, match map[string]interface{}) node.Match {
	return &node.TagMatch{
		Value:  unpack.OptionalString(ctx, match, key.Value, ""),
		Regex:  UnpackRegex(ctx, match),
		Prefix: unpack.OptionalString(ctx, match, key.Prefix, ""),
		Suffix: unpack.OptionalString(ctx, match, key.Suffix, ""),
	}
}

func UnpackMatchAttr(ctx *errctx.Context, match map[string]interface{}) node.Match {
	attr, aerr := maputil.RequireString(match, key.Name)
	if aerr != nil {
		ctx.ErrorWithKey(aerr, key.Name)
		return nil
	}

	return &node.AttrMatch{
		Attribute: attr,
		Value:     unpack.OptionalString(ctx, match, key.Value, ""),
		Regex:     UnpackRegex(ctx, match),
		Prefix:    unpack.OptionalString(ctx, match, key.Prefix, ""),
		Suffix:    unpack.OptionalString(ctx, match, key.Suffix, ""),
	}
}

func UnpackMatchAllOf(ctx *errctx.Context, match map[string]interface{}) node.Match {
	matches := UnpackMatchArray(ctx, match)
	if len(matches) == 0 {
		return nil
	}
	return node.AllOf(matches)
}

func UnpackMatchAnyOf(ctx *errctx.Context, match map[string]interface{}) node.Match {
	matches := UnpackMatchArray(ctx, match)
	if len(matches) == 0 {
		return nil
	}
	return node.AnyOf(matches)
}

func UnpackMatchOneOf(ctx *errctx.Context, match map[string]interface{}) node.Match {
	matches := UnpackMatchArray(ctx, match)
	if len(matches) == 0 {
		return nil
	}
	return node.OneOf(matches)
}

func UnpackMatchNot(ctx *errctx.Context, match map[string]interface{}) node.Match {
	submatch := unpack.RequireObject(ctx, match, key.Match)
	if submatch == nil {
		return nil
	}
	ctx.Path.Add(mpath.Key(key.Match))
	m := UnpackMatch(ctx, submatch)
	ctx.Path.Pop()
	if m == nil {
		return nil
	}

	return node.Not{Child: m}
}
