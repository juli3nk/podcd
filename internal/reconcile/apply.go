package reconcile

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/juli3nk/podcd/internal/model"
	"github.com/juli3nk/podcd/internal/source"
)

func (r *Reconciler) apply(spec model.ContainerSpec) error {
	// 1. init tasks
	for _, task := range spec.Init {
		if r.init.AlreadyDone(task) {
			continue
		}

		log.Printf("running init task: %s", task.Name)

		if err := r.init.Run(task); err != nil {
			return err
		}
	}

	// 2. render
	units, err := r.renderer.Render(spec)
	if err != nil {
		return err
	}

	// 3. write files
	for _, u := range units {
		if err := os.WriteFile(u.Path, []byte(u.Content), 0644); err != nil {
			return err
		}
	}

	// 4. reload systemd
	if err := r.systemd.Reload(); err != nil {
		return err
	}

	// 5. enable + start
	for _, u := range units {
		log.Printf("enabling service: %s", u.Name)

		if err := r.systemd.EnableNow(u.Name); err != nil {
			return err
		}
	}

	return nil
}

func (r *Reconciler) applyContainerChange(change source.Change) error {

	spec, err := loadSpecFromFile(r.source.Path(), change.Path)
	if err != nil {
		return err
	}

	switch change.Type {

	case source.Added, source.Modified:
		return r.applySpec(spec)

	case source.Deleted:
		return r.deleteSpec(spec)

	}

	return nil
}

func (r *Reconciler) applySpec(spec model.ContainerSpec) error {

	log.Printf("applying %s", spec.Name)

	// 1. INIT TASKS
	for _, task := range spec.Init {
		if r.init.AlreadyDone(task) {
			continue
		}

		log.Printf("running init task: %s", task.Name)

		if err := r.init.Run(task); err != nil {
			return err
		}
	}

	// 2. RENDER
	units, err := r.renderer.Render(spec)
	if err != nil {
		return err
	}

	// 3. WRITE FILES (avec diff)
	for _, u := range units {

		changed, err := writeIfChanged(u.Path, u.Content)
		if err != nil {
			return err
		}

		if !changed {
			log.Printf("unit unchanged: %s", u.Name)
			continue
		}

		log.Printf("updated unit: %s", u.Name)

		// 4. Enable + restart seulement si changé
		if err := r.systemd.Enable(true, u.Name); err != nil {
			return err
		}
	}

	return nil
}

func writeIfChanged(path, content string) (bool, error) {

	existing, err := os.ReadFile(path)
	if err == nil {
		if string(existing) == content {
			return false, nil
		}
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return false, err
	}

	return true, nil
}

func (r *Reconciler) deleteSpec(spec model.ContainerSpec) error {

	name := spec.Name + ".service"

	log.Printf("deleting %s", name)

	if err := r.systemd.Disable(name); err != nil {
		return err
	}

	if err := r.systemd.Stop(name); err != nil {
		return err
	}

	path := "/etc/systemd/system/" + name
	return os.Remove(path)
}

func filterContainerChanges(changes []source.Change) []source.Change {
	var out []source.Change

	for _, c := range changes {
		if strings.HasPrefix(c.Path, "containers/") {
			out = append(out, c)
		}
	}

	return out
}

func loadSpecFromFile(basePath, relPath string) (model.ContainerSpec, error) {

	fullPath := filepath.Join(basePath, relPath)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return model.ContainerSpec{}, err
	}

	var spec model.ContainerSpec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return model.ContainerSpec{}, err
	}

	return spec, nil
}
