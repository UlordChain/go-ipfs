package ca

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"net"

	"io"

	"encoding/base64"

	"github.com/pkg/errors"
)

const (
	requestMsgSize      = 85
	requestMsgVersion   = 111
	requestMsgQuestType = 1
	requestMsgTimestamp = 0
)

func writeString(w io.Writer, data string) {
	WriteVlen(w, uint64(len(data)))
	w.Write([]byte(data))
}

func MakeNodeInfoHash(txid string, voutid int32, pubkey string, licperiod int64, licversion int32) string {
	b := bytes.NewBuffer(nil)
	writeString(b, "Ulord Signed Message:\n")

	writeString(b, txid)

	binary.Write(b, binary.LittleEndian, voutid)

	_pubkey, _ := hex.DecodeString(pubkey)
	writeString(b, string(_pubkey))

	binary.Write(b, binary.LittleEndian, licperiod)

	binary.Write(b, binary.LittleEndian, licversion)

	uint256 := NewSha2Hash(b.Bytes())
	return uint256.String()
}

func MakeRandomHash() string {
	buf := make([]byte, 64)
	rand.Read(buf)
	return NewSha2Hash(buf).String()
}

func MakePrivateAddr(compress ...bool) string {
	key := make([]byte, 32)
	rand.Read(key)

	compre := false
	if len(compress) > 0 {
		compre = compress[0]
	}
	addr := NewPrivateAddr(key, 128, compre)
	//fmt.Println(
	return addr.String()
}

func PublicKeyFromPrivateAddr(privateAddr string) (string, error) {
	addr, err := DecodePrivateAddr(privateAddr)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(addr.Pubkey), nil
}

func VerifySignature(hash, licenseInBase64, pubkeyInHex string) (bool, error) {

	_, sign, err := ParseMessageSignature(licenseInBase64)
	if err != nil {
		return false, err
	}

	pubkeyBytes, err := hex.DecodeString(pubkeyInHex)
	if err != nil {
		return false, err
	}
	pubkey, err := NewPublicKey(pubkeyBytes)
	if err != nil {
		return false, err
	}

	return EcdsaVerify(pubkey.Bytes(true), sign.Bytes(), NewUint256FromString(hash).Bytes()), nil
}

type LicenseMetaInfo struct {
	MsgVersion int32 // 111
	Num        int32 // 1
	NodeType   int32 // 1
	Version    int32

	Txid       string // use
	Voutid     uint32 // use
	Privkey    string
	Status     int32
	Licversion int32  // use
	LicPeriod  int64  // use
	License    string // use
	Nodeperiod int64
}

func RequestLicense(servAddr string, txid string, voutid int32) (info *LicenseMetaInfo, e error) {
	type RequestMsg struct {
		size      int32
		version   int32
		timestamp int64
		questtype int32
		txid      string
		vountid   int32
	}

	msg := RequestMsg{
		size:      requestMsgSize,
		version:   requestMsgVersion,
		timestamp: requestMsgTimestamp,
		questtype: requestMsgQuestType,
		txid:      txid,
		vountid:   voutid,
	}

	b := bytes.NewBuffer(nil)
	binary.Write(b, binary.BigEndian, msg.size)
	binary.Write(b, binary.LittleEndian, msg.version)
	binary.Write(b, binary.LittleEndian, msg.timestamp)
	binary.Write(b, binary.LittleEndian, msg.questtype)
	WriteVlen(b, uint64(len(msg.txid)))
	b.Write([]byte(msg.txid))
	binary.Write(b, binary.LittleEndian, msg.vountid)

	if b.Len()-4 != int(msg.size) {
		e = errors.New("error request license message length")
		return
	}

	conn, err := net.Dial("tcp", servAddr)
	if err != nil {
		e = errors.Wrap(err, "dial license center service failed")
		return
	}
	defer conn.Close()

	_, err = conn.Write(b.Bytes())
	if err != nil {
		e = errors.Wrap(err, "write request license message failed")
		return
	}

	// read response
	var size int32
	err = binary.Read(conn, binary.BigEndian, &size)
	if err != nil {
		e = errors.Wrap(err, "read response license message failed")
		return
	}

	res := &LicenseMetaInfo{}

	err = binary.Read(conn, binary.LittleEndian, &res.MsgVersion)
	if err != nil {
		e = errors.Wrap(err, "read response license message field <msgversion> failed")
		return
	}

	err = binary.Read(conn, binary.LittleEndian, &res.Num)
	if err != nil {
		e = errors.Wrap(err, "read response license message field <num> failed")
		return
	}
	err = binary.Read(conn, binary.LittleEndian, &res.NodeType)
	if err != nil {
		e = errors.Wrap(err, "read response license message field <nodetype> failed")
		return
	}
	err = binary.Read(conn, binary.LittleEndian, &res.Version)
	if err != nil {
		e = errors.Wrap(err, "read response license message field <version> failed")
		return
	}

	res.Txid, err = ReadString(conn)
	if err != nil {
		e = errors.Wrap(err, "read response license message field <txid> failed")
		return
	}

	err = binary.Read(conn, binary.LittleEndian, &res.Voutid)
	if err != nil {
		e = errors.Wrap(err, "read response license message field <voutid> failed")
		return
	}

	res.Privkey, err = ReadString(conn)
	if err != nil {
		e = errors.Wrap(err, "read response license message field <privkey> failed")
		return
	}
	err = binary.Read(conn, binary.LittleEndian, &res.Status)
	if err != nil {
		e = errors.Wrap(err, "read response license message field <status> failed")
		return
	}
	err = binary.Read(conn, binary.LittleEndian, &res.Licversion)
	if err != nil {
		e = errors.Wrap(err, "read response license message field <licversion> failed")
		return
	}

	err = binary.Read(conn, binary.LittleEndian, &res.LicPeriod)
	if err != nil {
		e = errors.Wrap(err, "read response license message field <licperiod> failed")
		return
	}

	res.License, err = ReadString(conn)
	if err != nil {
		e = errors.Wrap(err, "read response license message field <license> failed")
		return
	}

	err = binary.Read(conn, binary.LittleEndian, &res.Nodeperiod)
	if err != nil {
		e = errors.Wrap(err, "read response license message field <nodeperiod> failed")
		return
	}

	info = res
	e = nil
	return
}

func Sign(hashInHex, pri string) (string, error) {
	signkey, err := DecodePrivateAddr(pri)
	if err != nil {
		return "", err
	}

	btcsig := new(Signature)
	var sb [65]byte
	sb[0] = 27
	if signkey.IsCompressed() {
		sb[0] += 4
	}
	r, s, err := EcdsaSign(signkey.Key, NewUint256FromString(hashInHex).Bytes())
	if err != nil {
		return "", err
	}
	btcsig.R.Set(r)
	btcsig.S.Set(s)

	rd := btcsig.R.Bytes()
	sd := btcsig.S.Bytes()
	copy(sb[1+32-len(rd):], rd)
	copy(sb[1+64-len(sd):], sd)

	return base64.StdEncoding.EncodeToString(sb[:]), nil
}

type UcenterPubkeyMap struct {
	MsgVersion int32 // 111
	Num        int32 //
	NodeType   int32 // 2

	V2key map[int32]string
}

func RequestUcenterPublicKeyMap(servAddr string, txid string, voutid int32) (info *UcenterPubkeyMap, e error) {
	type RequestMsg struct {
		size      int32
		version   int32
		timestamp int64
		questtype int32
		txid      string
		vountid   int32
	}

	msg := RequestMsg{
		size:      requestMsgSize,
		version:   requestMsgVersion,
		timestamp: requestMsgTimestamp,
		questtype: 2,
		txid:      txid,
		vountid:   voutid,
	}

	b := bytes.NewBuffer(nil)
	binary.Write(b, binary.BigEndian, msg.size)
	binary.Write(b, binary.LittleEndian, msg.version)
	binary.Write(b, binary.LittleEndian, msg.timestamp)
	binary.Write(b, binary.LittleEndian, msg.questtype)
	WriteVlen(b, uint64(len(msg.txid)))
	b.Write([]byte(msg.txid))
	binary.Write(b, binary.LittleEndian, msg.vountid)

	if b.Len()-4 != int(msg.size) {
		e = errors.New("error request license message length")
		return
	}

	conn, err := net.Dial("tcp", servAddr)
	if err != nil {
		e = errors.Wrap(err, "dial license center service failed")
		return
	}
	defer conn.Close()

	_, err = conn.Write(b.Bytes())
	if err != nil {
		e = errors.Wrap(err, "write request license message failed")
		return
	}

	// read response
	var size int32
	err = binary.Read(conn, binary.BigEndian, &size)
	if err != nil {
		e = errors.Wrap(err, "read response license message failed")
		return
	}

	res := &UcenterPubkeyMap{
		V2key: make(map[int32]string, 0),
	}

	err = binary.Read(conn, binary.LittleEndian, &res.MsgVersion)
	if err != nil {
		e = errors.Wrap(err, "read response license message field <msgversion> failed")
		return
	}

	err = binary.Read(conn, binary.LittleEndian, &res.Num)
	if err != nil {
		e = errors.Wrap(err, "read response license message field <num> failed")
		return
	}
	err = binary.Read(conn, binary.LittleEndian, &res.NodeType)
	if err != nil {
		e = errors.Wrap(err, "read response license message field <nodetype> failed")
		return
	}

	var (
		ver int32
		key string
	)
	for i := int32(0); i < res.Num; i++ {
		err = binary.Read(conn, binary.LittleEndian, &ver)
		if err != nil {
			e = errors.Wrapf(err, "read version failed from the %d index", i)
			return
		}

		key, err = ReadString(conn)
		if err != nil {
			e = errors.Wrapf(err, "read pubkey failed from the %s jindex", i)
			return
		}

		res.V2key[ver] = key
	}

	info = res
	e = nil
	return
}
