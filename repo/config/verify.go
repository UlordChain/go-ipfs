package config

type VerifyInfo struct {
	Txid   string
	Voutid int32
	Secret string

	Licversion int32
	License    string
	Period     int64
}
