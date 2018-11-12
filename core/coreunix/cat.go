package coreunix

import (
	"context"

	core "github.com/udfs/go-udfs/core"
	path "github.com/udfs/go-udfs/path"
	resolver "github.com/udfs/go-udfs/path/resolver"
	uio "github.com/udfs/go-udfs/unixfs/io"
)

func Cat(ctx context.Context, n *core.IpfsNode, pstr string) (uio.DagReader, error) {
	r := &resolver.Resolver{
		DAG:         n.DAG,
		ResolveOnce: uio.ResolveUnixfsOnce,
	}

	dagNode, err := core.Resolve(ctx, n.Namesys, r, path.Path(pstr))
	if err != nil {
		return nil, err
	}

	return uio.NewDagReader(ctx, dagNode, n.DAG)
}
