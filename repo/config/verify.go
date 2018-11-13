package config

type VerifyInfo struct {
	Txid   string
	Voutid int32
	Secret string

	Licversion int32  `json:",omitempty"`
	License    string `json:",omitempty"`
	Period     int64  `json:",omitempty"`
}
