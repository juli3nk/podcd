package reconcile

import (
	"log"
	"strings"
)

func (r *Reconciler) Reconcile() error {
	// 1. fetch Git
	if err := r.source.Fetch(); err != nil {
		return err
	}

	desired := loadFromGit()
    actual  := discoverFromRuntime()

    for name := range union(actual, desired) {

        desiredHash := hash(desired[name])
        actualHash  := actual[name].Hash

        switch {
        case not exists in actual:
            create()

        case not exists in desired:
            delete()

        case hash mismatch:
            recreate()

        case equal:
            noop()
        }
    }

	return nil
}
