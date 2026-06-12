package model

import "time"

type ContainerState struct {
	Name string

	// hash du spec appliqué
	LastAppliedHash string

	// statut runtime
	Running bool

	// timestamp utile
	LastUpdated time.Time
}
