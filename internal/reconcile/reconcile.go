package reconcile

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

	// 3. Discover actual state from runtime
	actual, err := r.discoverFromRuntime()
	if err != nil {
		return err
	}

	if err := r.networks.Reconcile(desired.Networks, actual.Networks); err != nil {
		return err
	}

	if err := r.volumes.Reconcile(desired.Volumes, actual.Volumes); err != nil {
		return err
	}

	if err := r.configMaps.Reconcile(desired.ConfigMaps, actual.ConfigMaps); err != nil {
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
