package reconcile

import (
	"os"
	"path/filepath"

	"github.com/juli3nk/podcd/internal/model"
	"github.com/juli3nk/podcd/internal/normalize"
)

func (r *ContainerReconciler) Reconcile(
	desired []model.Container,
	actual map[string]RuntimeObject,
) error {
	desiredMap := make(map[string]model.Container)
	for _, c := range desired {
		desiredMap[c.Name] = c
	}

	names := unionKeysContainer(actual, desiredMap)

	for _, name := range names {
		d, dExists := desiredMap[name]
		a, aExists := actual[name]

		desiredHash := normalize.HashContainer(d)

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
	serviceName := spec.Name + ".service"
	unitFullPath := filepath.Join(r.systemd.UnitPath(), serviceName)

	_ = r.systemd.Stop(serviceName)

	units, err := r.renderer.Render(spec, hash)
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
