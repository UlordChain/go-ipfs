package config

type VerifyInfo struct {
	ServerAddress string
	ServerPubkey  string

	Txid   string
	Voutid int32
	Secret string

	Licversion int32
	License    string
	Period     int64
}
