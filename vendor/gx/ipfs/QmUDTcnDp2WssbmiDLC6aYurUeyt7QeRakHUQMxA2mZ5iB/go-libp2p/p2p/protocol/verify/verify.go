package verify

import (
	"context"
	"time"

	"github.com/ipfs/go-ipfs/repo"
	"github.com/ipfs/go-ipfs/udfs/ca"
	"github.com/pkg/errors"

	config "gx/ipfs/QmPEpj17FDRpc7K1aArKZp3RsHtzRMKykeK9GVgn4WQGPR/go-ipfs-config"
	pb "gx/ipfs/QmUDTcnDp2WssbmiDLC6aYurUeyt7QeRakHUQMxA2mZ5iB/go-libp2p/p2p/protocol/verify/pb"
	inet "gx/ipfs/QmXuRkCR7BNQa9uqfpTiFWsTQLzmTWYg91Ja1w95gnqb6u/go-libp2p-net"
	logging "gx/ipfs/QmZChCsSt8DctjceaL56Eibc29CVQq4dGKRXC5JRZ6Ppae/go-log"
	msmux "gx/ipfs/QmabLh8TrJ3emfAoQk5AbqbLTbMyj7XqumMFmAFxa9epo8/go-multistream"
	host "gx/ipfs/QmdJfsSbKSZnMkfZ1kpopiyB9i3Hd6cp8VKWZmtWPa7Moc/go-libp2p-host"
	ggio "gx/ipfs/QmdxUuburamoF6zF9qjeQC4WYcWGbWuRmdLacMEsW8ioD8/gogo-protobuf/io"
)

const ProtocolVerify = "/ipfs/verify/0.0.1"

var log = logging.Logger("net/verify")

var (
	gcfg   *config.Config
	pubkey string
)

func readVerifyAskValue(s inet.Stream) (string, error) {
	r := ggio.NewDelimitedReader(s, 1024)
	mes := pb.VerifyAsk{}
	if err := r.ReadMsg(&mes); err != nil {
		return "", errors.Wrap(err, "read verify ask message failed")
	}

	return mes.GetRandom(), nil
}

func sendVerifyAskMsg(s inet.Stream) (string, error) {

	rhash := ca.MakeRandomHash()

	w := ggio.NewDelimitedWriter(s)
	mes := pb.VerifyAsk{
		Random: &rhash,
	}
	if err := w.WriteMsg(&mes); err != nil {
		return "", errors.Wrap(err, "write verify ask message failed")
	}

	return rhash, nil
}

func doVerify(mes *pb.Verify, rhash string, srvPubkey string) error {
	if time.Now().After(time.Unix(mes.GetPeriod(), 0)) {
		return errors.New("the license have out of date")
	}

	// verify license
	nodehash := ca.MakeNodeInfoHash(mes.GetTxid(), mes.GetVoutid(), mes.GetPubkey(), mes.GetPeriod(), mes.GetLicversion())
	if nodehash == "" {
		return errors.New("make node hash failed")
	}
	log.Debugf("got nodehash=<%s> with txid=<%s>, voutid=<%d>, pubkey=<%s>,  period=<%d>, licversion=<%d>\n",
		nodehash, mes.GetTxid(), mes.GetVoutid(), mes.GetPubkey(), mes.GetPeriod(), mes.GetLicversion())

	ok, err := ca.VerifySignature(nodehash, mes.GetLicense(), srvPubkey)
	if err != nil {
		return errors.Wrap(err, "verify license error")
	}

	if !ok {
		return errors.Errorf("verify license failed with nodehash=<%s>, license=<%s>, srvPubkey=<%s>\n", nodehash, mes.GetLicense(), srvPubkey)
	}

	// verify node signature
	ok, err = ca.VerifySignature(rhash, mes.GetRandomSign(), mes.GetPubkey())
	if err != nil {
		return errors.Wrap(err, "verify node signature error")
	}

	if !ok {
		return errors.Errorf("verify node signature failed with rhash=<%s>, rsigh=<%s>, pubkey=<%s>",
			rhash, mes.GetRandomSign(), mes.GetPubkey())
	}

	return nil
}

func getServerPubkey(c inet.Conn, h host.Host, licversion int32) (string, error) {
	for _, sp := range gcfg.UCenter.ServerPubkeys {
		if sp.Licversion == licversion {
			return sp.Pubkey, nil
		}
	}

	// request
	upm, err := ca.RequestUcenterPublicKeyMap(gcfg.UCenter.ServerAddress, gcfg.Verify.Txid, gcfg.Verify.Voutid)
	if err != nil {
		return "", errors.Wrap(err, "request ucenter public key map failed")
	}

	if len(upm.V2key) < len(gcfg.UCenter.ServerPubkeys) {
		return "", errors.New("request ucenter public key map less than already known")
	}

	var vps []*config.VersionPubkey
	var ret string
	for ver, key := range upm.V2key {
		vps = append(vps, &config.VersionPubkey{
			Licversion: ver,
			Pubkey:     key,
		})

		if ver == licversion {
			ret = key
		}
	}
	gcfg.UCenter.ServerPubkeys = vps

	// save
	repoInf, err := h.Peerstore().Get(c.LocalPeer(), "repo")
	if err != nil {
		return "", errors.Wrap(err, "got repo object from peerstore failed")
	}

	err = repoInf.(repo.Repo).SetConfig(gcfg)
	if err != nil {
		return "", errors.Wrap(err, "save config for new ucenter public key to file failed")
	}

	if ret == "" {
		return "", errors.Errorf("can`t get the ucenter public key for licversion=%d\n", licversion)
	}

	return ret, nil
}

// VerifyConn verify whether this connection is legal, if not, must close this connection
func VerifyConn(c inet.Conn, h host.Host) error {
	var err error

	// get the verify info
	if gcfg == nil {
		repoInf, err := h.Peerstore().Get(c.LocalPeer(), "repo")
		if err != nil {
			return errors.Wrap(err, "got repo object from peerstore failed")
		}
		repo := repoInf.(repo.Repo)
		cfg, err := repo.Config()
		if err != nil {
			return errors.Wrap(err, "got config object from repo failed")
		}
		gcfg = cfg
	}

	// get pubkey
	if pubkey == "" {
		pubkey, err = ca.PublicKeyFromPrivateAddr(gcfg.Verify.Secret)
		if err != nil {
			return errors.Wrap(err, "got public key from secret failed")
		}
	}

	// make new stream
	s, err := c.NewStream()
	if err != nil {
		log.Debugf("error opening initial stream for %s: %s", ProtocolVerify, err)
		log.Event(context.TODO(), "VerifyConnFailed", c.RemotePeer())
		return errors.Wrap(err, "new stream failed")
	}
	defer inet.FullClose(s)

	s.SetProtocol(ProtocolVerify)

	err = s.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return errors.Wrap(err, "set read deadline failed")
	}
	err = s.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return errors.Wrap(err, "set write deadline failed")
	}

	// ok give the response to our handler.
	if err := msmux.SelectProtoOrFail(ProtocolVerify, s); err != nil {
		log.Event(context.TODO(), "VerifyConnFailed", c.RemotePeer(), logging.Metadata{"error": err})
		return errors.Wrap(err, "SelectProtoOrFail for verify failed")
	}

	// send verify ask msg
	rhash, err := sendVerifyAskMsg(s)
	if err != nil {
		return errors.Wrap(err, "send verify ask msg failed")
	}

	// read verify info
	r := ggio.NewDelimitedReader(s, 2048)
	mes := &pb.Verify{}
	if err := r.ReadMsg(mes); err != nil {
		return errors.Wrap(err, "error reading verify message")
	}

	log.Debugf("%s received message from %s %s", ProtocolVerify,
		c.RemotePeer(), c.RemoteMultiaddr())

	log.Debugf("received verify msg for %s : %v\n", c.RemotePeer(), mes)

	ucenterPubkey, err := getServerPubkey(c, h, mes.GetLicversion())
	if err != nil {
		return errors.Wrap(err, "get ucenter public key failed")
	}
	err = doVerify(mes, rhash, ucenterPubkey)
	if err != nil {
		return errors.Wrapf(err, "verify %v failed", c.RemotePeer())
	}

	return nil
}

func CheckVerifyInfo(vfi *config.VerifyInfo) error {
	if len(vfi.Txid) == 0 {
		return errors.New("the field <Txid> in verify info is not set")
	}

	if len(vfi.Secret) == 0 {
		return errors.New("the field <Secret> in verify info is not set")
	}

	if vfi.Voutid < 0 {
		return errors.New("the field <Voutid> in verify info is not legal")
	}

	return nil
}

func CheckUCenterInfo(info *config.UCenterInfo) error {
	if info.ServerAddress == "" {
		return errors.New("the field <ServerAddress> in verify into is not set")
	}

	return nil
}

func RequestHandler(s inet.Stream, h host.Host) {
	defer s.Close()

	// read verify ask msg
	rhash, err := readVerifyAskValue(s)
	if err != nil {
		log.Error("read verify ask value failed: ", err)
		return
	}

	err = requestVerify(s, h, rhash)
	if err != nil {
		log.Error("request verify failed: ", err)
		return
	}
}

func requestVerify(s inet.Stream, h host.Host, rhash string) error {
	// get the verify info
	vfi := &gcfg.Verify

	// request license
	var err error

	// check need request license
	if time.Now().After(time.Unix(vfi.Period, 0)) || vfi.License == "" || vfi.Licversion == 0 {
		log.Debug("request license...")

		lbi, err := ca.RequestLicense(gcfg.UCenter.ServerAddress, vfi.Txid, vfi.Voutid)
		if err != nil {
			return errors.Wrap(err, "request license failed")
		}
		vfi.License = lbi.License
		vfi.Period = lbi.LicPeriod
		vfi.Licversion = lbi.Licversion

		repoInf, err := h.Peerstore().Get(s.Conn().LocalPeer(), "repo")
		if err != nil {
			return errors.Wrap(err, "got repo object from peerstore failed")
		}
		err = repoInf.(repo.Repo).SetConfig(gcfg)
		if err != nil {
			return errors.Wrap(err, "save config for new license to file failed")
		}
	}

	// make signature for random hash
	sign, err := ca.Sign(rhash, vfi.Secret)
	if err != nil {
		return errors.Wrap(err, "sign for rhash failed")
	}

	// send verify msg
	w := ggio.NewDelimitedWriter(s)
	mes := pb.Verify{
		Txid:       &vfi.Txid,
		Voutid:     &vfi.Voutid,
		Period:     &vfi.Period,
		Pubkey:     &pubkey,
		License:    &vfi.License,
		Licversion: &vfi.Licversion,
		RandomSign: &sign,
	}
	w.WriteMsg(&mes)

	return nil
}
