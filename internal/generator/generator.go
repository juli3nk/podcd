package generator

import "github.com/juli3nk/podcd/internal/model"

type Generator interface {
	Generate(spec model.SecretGenerator) ([]byte, error)
}
