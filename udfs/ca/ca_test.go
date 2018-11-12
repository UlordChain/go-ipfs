package ca

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"testing"
)

var (
	txid   = "829f9b963be709e8b138385fedb876990f892b99ca641401e8a88df531735ca7"
	voutid = int32(1)

	period     = int64(1543386972)
	licversion = int32(1)

	pubkeyStr = "0427989b89ebc4c8596a82830a877b5edbde0739ea8df6e7c8bbbdf3bf0857f555bb739f3e9ce9951a0f34be4ff4d7e771b544401139de4dcca4809a0aac4c6ed5"
	hashStr   = "6b1b03d59c397fa77967a10e33259c58db295c1f8b66ba568a6c13ea0cc8ea21"

	signStr = "IEC1TEc+iLcJ2LsXFZ2Xdx7E8VmXuoW98unHqsh2uNXmBmMpI7Lm+/v1B7hrczJ0hO3pMh5yTgk/6aS7L/VnCGE="

	ucenterPubkeyStr = "03e947099921ee170da47a7acf48143c624d33950af362fc39a734b1b3188ec1e3"
)

func Test_MakeBTCHash(t *testing.T) {
	if MakeNodeInfoHash(txid, voutid, pubkeyStr, period, licversion) != hashStr {
		t.Failed()
	}
}

func Test_RequestSignature(t *testing.T) {
	p, lic, e := RequestLicense(txid, voutid)
	if e != nil {
		t.Fatal(e)
	}

	if p != period {
		t.Failed()
	}

	if lic != signStr {
		t.Failed()
	}

}

func Test_VerifySignature(t *testing.T) {
	ok, err := VerifySignature(hashStr, signStr, ucenterPubkeyStr)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Failed()
	}
}

func Test_Raw(t *testing.T) {
	b := bytes.NewBuffer(nil)
	// strMessageMagic
	strMessageMagic := "Ulord Signed Message:\n"
	WriteVlen(b, uint64(len(strMessageMagic)))
	b.Write([]byte(strMessageMagic))

	// _txid
	txid := "da29e46c4d805474cbed1e6771473c77b350858e990de26837fcbeb8a0271acf"
	WriteVlen(b, uint64(len(txid)))
	b.Write([]byte(txid))

	// _voutid
	voutid := uint32(0)
	binary.Write(b, binary.LittleEndian, uint32(voutid))

	// _pubkey
	pubkey := "041df129b64bf4091c6093b7c6253479b96541ce1fef817a9b0936754b7f25d6f43979666350621e0a59541868916f1b9b39de0e92dc2d108aec1177c9baea7858"
	_pubkey, _ := hex.DecodeString(pubkey)
	WriteVlen(b, uint64(len(_pubkey)))
	b.Write(_pubkey)

	// _licperiod
	licperiod := int64(1541648532)
	binary.Write(b, binary.LittleEndian, int64(licperiod))

	// _licversion
	licversion := int32(1)
	binary.Write(b, binary.LittleEndian, int32(licversion))

	// fmt.Println(b.Len())
	uint256 := NewSha2Hash(b.Bytes())
	fmt.Println(uint256.String())

	if uint256.String() != "cd77355ebc70f9d1b674a4d3dacfb09464dcfbf2b3055ab97ee8a7d28e447784" {
		fmt.Println("failed")
	}

	// cd77355ebc70f9d1b674a4d3dacfb09464dcfbf2b3055ab97ee8a7d28e447784
}

func Test_MakePrivateAddr(t *testing.T) {
	fmt.Println(MakePrivateAddr())
}

func Test_PublicKeyFromPrivateAddr(t *testing.T) {

	pri := ("5KR4UPG5GbMakkkrsHuzmPZy9bLmkXAth56r99cZAvRRkNk5581")
	pub := PublicKeyFromPrivateAddr(string(pri))
	if pub != "042660fa685a470ff43ee6dacd6b8979ade29f8ce2bf1fd8f6844ed29983153fb8e5f3e90af800e892946dae75411dd1acf58559afc22166fb13ca21696f8880e7" {
		t.FailNow()
	}
}

func Test_Sign(t *testing.T) {
	pri := ("5KR4UPG5GbMakkkrsHuzmPZy9bLmkXAth56r99cZAvRRkNk5581")
	pub := PublicKeyFromPrivateAddr(string(pri))

	hash := MakeRandomHash()
	fmt.Println(len(hash), hash)
	b64, err := Sign(hash, pri)
	if err != nil {
		t.Fatal(err)
	}

	ok, err := VerifySignature(hash, b64, pub)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.FailNow()
	}
}
