package model

type Secret struct {
	Name string

	EncryptedData map[string]string

	Labels   map[string]string
	Metadata Metadata
}
