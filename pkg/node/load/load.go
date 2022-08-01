package load

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/tvarney/maputil/consterr"
	"github.com/tvarney/maputil/errctx"
	"github.com/tvarney/sdtdmod/pkg/node"
	"github.com/tvarney/sdtdmod/pkg/node/load/impl"
)

func LoadFile(filename string, handlers ...errctx.ErrorHandler) ([]*node.Node, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return LoadBytes(filename, data, handlers...)
}

func LoadReader(name string, r io.Reader, handlers ...errctx.ErrorHandler) ([]*node.Node, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return LoadBytes(name, data, handlers...)
}

func LoadValue(name string, value interface{}, handlers ...errctx.ErrorHandler) ([]*node.Node, error) {
	ctx := impl.CreateErrCtx(name, handlers...)
	switch v := value.(type) {
	case []interface{}:
		nl := impl.UnpackNodeList(ctx, v)
		return nl, impl.GetError(ctx)
	case map[string]interface{}:
		n := impl.UnpackNode(ctx, v)
		if n == nil {
			return nil, impl.GetError(ctx)
		}
		return []*node.Node{n}, impl.GetError(ctx)
	}
	return nil, consterr.Error("")
}

func LoadBytes(name string, data []byte, handlers ...errctx.ErrorHandler) ([]*node.Node, error) {
	ctx := impl.CreateErrCtx(name, handlers...)
	var arr []interface{}
	if err := json.Unmarshal(data, &arr); err == nil {
		nl := impl.UnpackNodeList(ctx, arr)
		return nl, impl.GetError(ctx)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	n := impl.UnpackNode(ctx, m)
	if n != nil {
		return []*node.Node{n}, impl.GetError(ctx)
	}
	return nil, impl.GetError(ctx)
}
