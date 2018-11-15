/*
Package coreapi provides direct access to the core commands in UDFS. If you are
embedding UDFS directly in your Go program, this package is the public
interface you should use to read and write files or otherwise control UDFS.

If you are running UDFS as a separate process, you should use `go-udfs-api` to
work with it via HTTP. As we finalize the interfaces here, `go-udfs-api` will
transparently adopt them so you can use the same code with either package.

**NOTE: this package is experimental.** `go-udfs` has mainly been developed
as a standalone application and library-style use of this package is still new.
Interfaces here aren't yet completely stable.
*/
package coreapi

import (
	core "github.com/udfs/go-udfs/core"
	coreiface "github.com/udfs/go-udfs/core/coreapi/interface"
)

type CoreAPI struct {
	node *core.UdfsNode
}

// NewCoreAPI creates new instance of UDFS CoreAPI backed by go-udfs Node.
func NewCoreAPI(n *core.UdfsNode) coreiface.CoreAPI {
	api := &CoreAPI{n}
	return api
}

// Unixfs returns the UnixfsAPI interface implementation backed by the go-udfs node
func (api *CoreAPI) Unixfs() coreiface.UnixfsAPI {
	return (*UnixfsAPI)(api)
}

// Block returns the BlockAPI interface implementation backed by the go-udfs node
func (api *CoreAPI) Block() coreiface.BlockAPI {
	return (*BlockAPI)(api)
}

// Dag returns the DagAPI interface implementation backed by the go-udfs node
func (api *CoreAPI) Dag() coreiface.DagAPI {
	return (*DagAPI)(api)
}

// Name returns the NameAPI interface implementation backed by the go-udfs node
func (api *CoreAPI) Name() coreiface.NameAPI {
	return (*NameAPI)(api)
}

// Key returns the KeyAPI interface implementation backed by the go-udfs node
func (api *CoreAPI) Key() coreiface.KeyAPI {
	return (*KeyAPI)(api)
}

// Object returns the ObjectAPI interface implementation backed by the go-udfs node
func (api *CoreAPI) Object() coreiface.ObjectAPI {
	return (*ObjectAPI)(api)
}

// Pin returns the PinAPI interface implementation backed by the go-udfs node
func (api *CoreAPI) Pin() coreiface.PinAPI {
	return (*PinAPI)(api)
}
