package secret

import "go.yaml.in/yaml/v4"

func LoadSecretData(path string, decrypter Decrypter) ([]byte, error) {
	data, err := decrypter.Decrypt(path)
	if err != nil {
		return nil, err
	}

	var payload struct {
		Data string `yaml:"data"`
	}

	if err := yaml.Unmarshal(data, &payload); err != nil {
		return nil, err
	}

	return []byte(payload.Data), nil
}
