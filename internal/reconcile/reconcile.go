package reconcile

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/juli3nk/podcd/internal/model"
	"github.com/juli3nk/podcd/internal/normalize"
	"github.com/juli3nk/podcd/internal/secret"
)

func (r *Reconciler) Reconcile() error {
	// 1. Fetch Git
	if err := r.source.Fetch(); err != nil {
		return err
	}

	// 2. Load desired state
	desired, err := loadFromGit(r.source.Path())
	if err != nil {
		return err
	}
	fmt.Printf("Desired: %+v\n", desired)

	// 3. Discover actual state from runtime
	actual, err := r.discoverFromRuntime()
	if err != nil {
		return err
	}
	fmt.Printf("Actual: %+v\n", actual)

	if err := r.networks.Reconcile(desired.Networks, actual.Networks); err != nil {
		return err
	}

	if err := r.volumes.Reconcile(desired.Volumes, actual.Volumes); err != nil {
		return err
	}

	if err := r.secrets.Reconcile(desired.Secrets, actual.Secrets); err != nil {
		return err
	}

	if err := r.containers.Reconcile(desired.Containers, actual.Containers); err != nil {
		return err
	}

	return nil
}

func (r *NetworkReconciler) Reconcile(
	desired []model.Network,
	actual map[string]RuntimeObject,
) error {
	desiredMap := make(map[string]model.Network)
	for _, net := range desired {
		desiredMap[net.Name] = net
	}

	for _, name := range unionKeysNetwork(actual, desiredMap) {
		desiredNet, desiredExists := desiredMap[name]
		actualNet, actualExists := actual[name]

		desiredHash := normalize.HashNetwork(desiredNet)

		switch {
		case !actualExists && desiredExists:
			if err := r.runtime.CreateNetwork(desiredNet, desiredHash); err != nil {
				return err
			}

		case actualExists && !desiredExists:
			if err := r.runtime.RemoveNetwork(name); err != nil {
				return err
			}

		case actualExists && desiredExists:
			// UPDATE (rare pour network → souvent recreate)
			if desiredHash != actualNet.Hash {
				if err := r.runtime.RemoveNetwork(name); err != nil {
					return err
				}
				if err := r.runtime.CreateNetwork(desiredNet, desiredHash); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (r *VolumeReconciler) Reconcile(
	desired []model.Volume,
	actual map[string]RuntimeObject,
) error {
	desiredMap := make(map[string]model.Volume)
	for _, v := range desired {
		desiredMap[v.Name] = v
	}

	names := unionKeysVolume(actual, desiredMap)

	for _, name := range names {
		d, dExists := desiredMap[name]
		a, aExists := actual[name]

		desiredHash := normalize.HashVolume(d)

		switch {
		case !aExists && dExists:
			if err := r.runtime.CreateVolume(d, desiredHash); err != nil {
				return err
			}

		case aExists && !dExists:
			if err := r.runtime.RemoveVolume(name); err != nil {
				return err
			}

		case aExists && dExists:
			// en pratique, volume = rarement modifiable
			if desiredHash != a.Hash {
				// ⚠️ attention destructive
				// → souvent mieux de skip ou warning
				if err := r.runtime.RemoveVolume(name); err != nil {
					return err
				}
				if err := r.runtime.CreateVolume(d, desiredHash); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

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

func (r *ContainerReconciler) Reconcile(
	desired []model.Container,
	actual map[string]RuntimeObject,
) error {
	desiredMap := make(map[string]model.Container)
	for _, c := range desired {
		desiredMap[c.Spec.Name] = c
	}

	names := unionKeysContainer(actual, desiredMap)

	for _, name := range names {
		d, dExists := desiredMap[name]
		a, aExists := actual[name]

		desiredHash := normalize.HashContainer(d.Spec)

		switch {
		case !aExists && dExists:
			if err := r.Create(d, desiredHash); err != nil {
				return err
			}

		case aExists && !dExists:
			if err := r.Delete(name); err != nil {
				return err
			}

		case aExists && dExists:
			if desiredHash != a.Hash {
				if err := r.Create(d, desiredHash); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (r *ContainerReconciler) Create(spec model.Container, hash string) error {
	serviceName := spec.Spec.Name + ".service"
	unitFullPath := filepath.Join(r.systemd.UnitPath(), serviceName)

	_ = r.systemd.Stop(serviceName)

	units, err := r.renderer.Render(spec.Spec, hash)
	if err != nil {
		return err
	}

	for _, u := range units {
		if err := os.WriteFile(unitFullPath, []byte(u.Content), 0644); err != nil {
			return err
		}
	}

	if err := r.systemd.Reload(); err != nil {
		return err
	}

	return r.systemd.Enable(serviceName)
}

func (r *ContainerReconciler) Delete(name string) error {
	serviceName := name + ".service"
	unitFullPath := filepath.Join(r.systemd.UnitPath(), serviceName)

	if err := r.systemd.Disable(serviceName); err != nil {
		return err
	}

	if err := r.systemd.Stop(serviceName); err != nil {
		return err
	}

	return os.Remove(unitFullPath)
}
