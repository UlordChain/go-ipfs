package ca

import (
	"fmt"
	"log"
	"testing"

	"time"

	"github.com/pkg/errors"
)

var (
	txid    = "ebfad1c5e579b42837625b421d58420f4d41cb0b2c488e74ad0bd9d5291948eb"
	voutid  = int32(0)
	privStr = "5JCRtN29x3QRxY5uyijQGbmyU7e1pN8EXWcifTsfjqKFnYUefyZ"

	testServerAddress = "132.232.98.139:5009"

	localtestSecret = "KzAzCLFC7g4sRqdvEdngGrG2YHJYJ81i4R6Y8Kr3Xg418pB8uhd1"
)

var secretMap = make(map[string]string, 0)

func init() {
	// register secret
	secretMap[txid] = privStr
}

func getUcenterPubkey(licversion int32) string {
	if licversion == 999 {
		p, _ := PublicKeyFromPrivateAddr(localtestSecret)
		return p
	}

	upm, err := RequestUcenterPublicKeyMap(testServerAddress, txid, voutid)
	if err != nil {
		log.Fatal(err)
	}

	pk, found := upm.V2key[licversion]
	if !found {
		log.Fatalf("can`t found pubkey for licversion=%d", licversion)
	}
	return pk
}

func mocRequestLicense(srvAddr, txid string, voutid int32) (info *LicenseMetaInfo, e error) {
	if srvAddr != "" {
		return RequestLicense(testServerAddress, txid, voutid)
	}

	secret, found := secretMap[txid]
	if !found {
		return nil, errors.Errorf("can`t found the secret for txid=%s", txid)
	}

	// got pubkey
	pubkeyStr, err := PublicKeyFromPrivateAddr(secret)
	if err != nil {
		return nil, errors.Wrap(err, "got pubkey failed")
	}

	// make node hash
	now := time.Now()
	year, month, day := now.Date()
	period := time.Date(year, month, day+1, 0, 0, 0, 0, now.Location()).Unix()
	nodeHash := MakeNodeInfoHash(txid, voutid, pubkeyStr, period, 999)

	//do sign
	s, err := Sign(nodeHash, localtestSecret)
	if err != nil {
		return nil, errors.Wrap(err, "sign failed")
	}

	lbi := &LicenseMetaInfo{
		License:    s,
		Licversion: 999,
		LicPeriod:  period,
		Txid:       txid,
		Voutid:     uint32(voutid),
	}

	return lbi, nil
}

func Test_Main(t *testing.T) {
	// request license
	lbi, e := mocRequestLicense(testServerAddress, txid, voutid)
	if e != nil {
		t.Fatal(e)
	}

	fmt.Println("txid:", txid)
	fmt.Println("voutid:", voutid)

	period := lbi.LicPeriod
	licversion := lbi.Licversion
	license := lbi.License

	fmt.Println("period:", period)
	fmt.Println("licversion:", licversion)
	fmt.Println("license:", license)

	secret, found := secretMap[txid]
	if !found {
		t.Fatal(errors.Errorf("can`t found the secret for txid=%s", txid))
	}
	fmt.Println("secret:", secret)

	// got pubkey
	pubkeyStr, err := PublicKeyFromPrivateAddr(secret)
	if err != nil {
		t.Fatal(errors.Wrap(err, "got pubkey failed"))
	}
	fmt.Println("pubkeyStr:", pubkeyStr)

	// make node hash
	nodeHash := MakeNodeInfoHash(txid, voutid, pubkeyStr, period, licversion)
	fmt.Println("nodeHash:", nodeHash)

	// verify license
	ok, err := VerifySignature(nodeHash, license, getUcenterPubkey(licversion))
	if err != nil {
		t.Fatal(errors.Wrap(err, "verify license error"))
	}

	if !ok {
		t.Fatal(errors.New("verify license failed"))
	}

	fmt.Println("finish")
}

func Test_MakePrivateAddr(t *testing.T) {
	fmt.Println(MakePrivateAddr(true))
}

func Test_RequestUcenterPublicKeyMap(t *testing.T) {
	fmt.Println(RequestUcenterPublicKeyMap(testServerAddress, txid, voutid))
}
