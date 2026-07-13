package model

type SecretGenerator struct {
	Type   string
	Length int

	Upper   bool
	Lower   bool
	Digits  bool
	Symbols bool
}

type Secret struct {
	Name string

	Filepath  string
	Generator *SecretGenerator

	Labels   map[string]string
	Metadata Metadata
}
