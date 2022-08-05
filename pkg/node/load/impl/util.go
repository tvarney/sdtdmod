package impl

import (
	"log"
	"regexp"

	"github.com/tvarney/maputil/errctx"
	"github.com/tvarney/maputil/mpath"
	"github.com/tvarney/maputil/unpack"
	"github.com/tvarney/sdtdmod/pkg/errors"
	"github.com/tvarney/sdtdmod/pkg/node"
	"github.com/tvarney/sdtdmod/pkg/node/key"
)

// CreateErrCtx creates an error context for loading nodes.
func CreateErrCtx(name string, handlers ...errctx.ErrorHandler) *errctx.Context {
	ctx := errctx.New(handlers...)
	ctx.Path.Filename = name
	if ctx.Handler == nil {
		log.Printf("Using default error collector")
		ctx.Handler = &errors.ErrorCollector{}
	}
	return ctx
}

// GetError returns the error that should be returned by load functions.
func GetError(ctx *errctx.Context) error {
	count := ctx.ErrorCount()
	if count == 0 {
		return nil
	}
	e := &errors.ParseCountError{
		Count: count,
		Name:  ctx.Path.Filename,
	}
	if h, ok := ctx.Handler.(*errors.ErrorCollector); ok {
		e.Errors = h.Errors
	}
	return e
}

// UnpackRegex unpacks a regex
func UnpackRegex(ctx *errctx.Context, obj map[string]interface{}) *regexp.Regexp {
	raw := unpack.OptionalString(ctx, obj, key.Regex, "")
	if raw == "" {
		return nil
	}
	r, err := regexp.Compile(raw)
	if err != nil {
		ctx.ErrorWithKey(err, key.Regex)
		return nil
	}
	return r
}

// UnpackCondition fetches the 'if' key and if present converts it to a Match.
func UnpackCondition(ctx *errctx.Context, action map[string]interface{}) node.Match {
	rawCond := unpack.OptionalObject(ctx, action, key.Cond, nil)
	if rawCond == nil {
		return nil
	}

	ctx.Path.Add(mpath.Key(key.Cond))
	cond := UnpackMatch(ctx, rawCond)
	ctx.Path.Pop()
	return cond
}
