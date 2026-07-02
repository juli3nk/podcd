package secret

type Decrypter interface {
	Decrypt(path string) ([]byte, error)
}
