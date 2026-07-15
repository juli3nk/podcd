package reconcile

import (
	"fmt"

	"github.com/juli3nk/podcd/internal/model"
	"github.com/juli3nk/podcd/internal/normalize"
	"github.com/juli3nk/podcd/internal/secret"
)

func (r *SecretReconciler) Reconcile(
	desired []model.Secret,
	actual map[string]RuntimeObject,
) error {
	desiredMap := make(map[string]model.Secret)
	for _, s := range desired {
		desiredMap[s.Name] = s
	}

	names := unionKeysSecret(actual, desiredMap)

	for _, name := range names {
		d, dExists := desiredMap[name]
		a, aExists := actual[name]

		secretData, err := r.secretData(d)
		if err != nil {
			return err
		}

		desiredHash := normalize.HashSecret(d)

		switch {
		case !aExists && dExists:
			if err := r.runtime.CreateSecret(d, secretData, desiredHash); err != nil {
				return err
			}

		case aExists && !dExists:
			if err := r.runtime.RemoveSecret(name); err != nil {
				return err
			}

		case aExists && dExists:
			// en pratique, volume = rarement modifiable
			if desiredHash != a.Hash {
				// ⚠️ attention destructive
				// → souvent mieux de skip ou warning
				if err := r.runtime.RemoveSecret(name); err != nil {
					return err
				}
				if err := r.runtime.CreateSecret(d, secretData, desiredHash); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (r *SecretReconciler) secretData(
	s model.Secret,
) ([]byte, error) {
	switch {
	case s.Filepath != "":
		return secret.LoadSecretData(
			s.Filepath,
			r.decrypter,
		)

	case s.Generator != nil:
		return r.generator.Generate(
			*s.Generator,
		)

	default:
		return nil, fmt.Errorf(
			"secret %q: neither filepath nor generator defined",
			s.Name,
		)
	}
}
