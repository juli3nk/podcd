package reconcile

import (
	"github.com/juli3nk/podcd/internal/model"
	"github.com/juli3nk/podcd/internal/normalize"
)

func (r *ConfigMapReconciler) Reconcile(
	desired []model.ConfigMap,
	actual map[string]RuntimeObject,
) error {
	desiredMap := make(map[string]model.ConfigMap)
	for _, cm := range desired {
		desiredMap[cm.Name] = cm
	}

	for _, name := range unionKeysConfigMap(actual, desiredMap) {
		desiredConfigMap, desiredExists := desiredMap[name]
		actualConfigMap, actualExists := actual[name]

		desiredHash := normalize.HashConfigMap(desiredConfigMap)

		switch {
		case !actualExists && desiredExists:
			if err := r.runtime.CreateConfigMap(desiredConfigMap, desiredHash); err != nil {
				return err
			}

		case actualExists && !desiredExists:
			if err := r.runtime.RemoveConfigMap(name); err != nil {
				return err
			}

		case actualExists && desiredExists:
			if desiredHash != actualConfigMap.Hash {
				if err := r.runtime.RemoveConfigMap(name); err != nil {
					return err
				}
				if err := r.runtime.CreateConfigMap(desiredConfigMap, desiredHash); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
