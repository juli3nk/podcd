package model

type IPAMConfig struct {
	Subnet  string
	Gateway string
	IPRange string
}

type Network struct {
	Name       string
	Driver     string
	IPAM       []IPAMConfig
	Internal   bool
	Attachable bool
	Labels     map[string]string
}
