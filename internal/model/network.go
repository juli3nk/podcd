package model

type IPAMConfig struct {
	Subnet  string
	Gateway string
	IPRange string
}

type IPAM struct {
	Driver  string
	Options map[string]string
	Config  []IPAMConfig
}

type Network struct {
	Name string

	Driver  string
	Options map[string]string

	IPAM IPAM

	Internal   bool
	Attachable bool

	Labels map[string]string
}
