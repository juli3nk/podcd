package generator

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/juli3nk/podcd/internal/model"
)

type RandomGenerator struct{}

func (g *RandomGenerator) Generate(
	spec model.SecretGenerator,
) ([]byte, error) {
	if spec.Length <= 0 {
		return nil, fmt.Errorf("length must be greater than 0")
	}

	var charset string

	if spec.Lower {
		charset += "abcdefghijklmnopqrstuvwxyz"
	}

	if spec.Upper {
		charset += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}

	if spec.Digits {
		charset += "0123456789"
	}

	if spec.Symbols {
		charset += "!@#$%^&*()-_=+[]{}<>?"
	}

	if charset == "" {
		return nil, fmt.Errorf("no character set selected")
	}

	result := make([]byte, spec.Length)

	for i := range result {
		n, err := rand.Int(
			rand.Reader,
			big.NewInt(int64(len(charset))),
		)
		if err != nil {
			return nil, err
		}

		result[i] = charset[n.Int64()]
	}

	return result, nil
}
