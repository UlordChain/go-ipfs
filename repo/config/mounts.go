package config

// Mounts stores the (string) mount points
type Mounts struct {
	UDFS           string
	IPNS           string
	FuseAllowOther bool
}
