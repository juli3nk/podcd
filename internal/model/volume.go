package model

type Volume struct {
	Name string

	Driver  string
	Options map[string]string

	Labels map[string]string
}
