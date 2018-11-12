package main

import (
	"context"
	"gx/ipfs/QmeiCcJfDW1GJnWUArudsv5rQsihpi4oyddPhdqo3CfX6i/go-datastore"
	syncds "gx/ipfs/QmeiCcJfDW1GJnWUArudsv5rQsihpi4oyddPhdqo3CfX6i/go-datastore/sync"
	"log"

	"fmt"

	"os"

	"github.com/udfs/go-udfs/core"
	"github.com/udfs/go-udfs/core/coreunix"
	"github.com/udfs/go-udfs/repo"
	"github.com/udfs/go-udfs/repo/config"
)

func main() {
	if len(os.Args) == 1 {
		log.Fatal("must provide file")
	}
	filename := os.Args[1]

	r := &repo.Mock{
		C: config.Config{
			Identity: config.Identity{
				PeerID: "QmTFauExutTsy4XP6JbMFcw2Wa9645HJt2bTqL6qYDCKfe", // required by offline node
			},
		},
		D: syncds.MutexWrap(datastore.NewMapDatastore()),
	}
	node, err := core.NewNode(context.Background(), &core.BuildCfg{Repo: r})
	if err != nil {
		log.Fatal(err)
	}
	k, err := coreunix.AddR(node, filename)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(k)
}
