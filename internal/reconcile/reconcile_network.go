package reconcile

import (
	"github.com/juli3nk/podcd/internal/model"
	"github.com/juli3nk/podcd/internal/normalize"
)

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
