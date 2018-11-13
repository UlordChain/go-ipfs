package config

type VersionPubkey struct {
	Licversion int32
	Pubkey     string
}
type UCenterInfo struct {
	ServerAddress string
	ServerPubkeys []*VersionPubkey
}
