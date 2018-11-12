PB = $(wildcard *.proto)
GO = $(PB:.proto=.pb.go)

all: $(GO)

%.pb.go: %.proto
		protoc --gogo_out=. --proto_path=../../../vendor/gx/ipfs/QmVsp2KdPYE6M8ryzCk5KHLo3zprcY5hBDaYx6uPCFUdxA/go-libp2p-record/pb:/usr/local/opt/protobuf/include:. $<

clean:
		rm -f *.pb.go