package model

type Secret struct {
	Name string

	Filepath string

	Labels   map[string]string
	Metadata Metadata
}
