package config

var UCenterServerAddress = "ucenter.ulord.one:5009"
type VersionPubkey struct {
	Licversion int32
	Pubkey     string
}
type UCenterInfo struct {
	ServerAddress string
	ServerPubkeys []*VersionPubkey
}
