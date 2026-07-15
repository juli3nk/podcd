package reconcile

import (
	"github.com/juli3nk/podcd/internal/model"
	"github.com/juli3nk/podcd/internal/normalize"
)

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
