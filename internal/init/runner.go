package initrunner

import (
	"github.com/juli3nk/podcd/internal/model"
	"github.com/juli3nk/podcd/internal/normalize"
	"github.com/juli3nk/podcd/internal/runtime"
	"github.com/juli3nk/podcd/internal/state"
)

type InitRunner interface {
	AlreadyDone(task model.InitTask) bool
	Run(task model.InitTask) error
}

type Runner struct {
	runtime runtime.Runtime
	state   state.Store
}

func New(rt runtime.Backend, state state.Store) (InitRunner, error) {
	rtImpl, err := runtime.New(rt)
	if err != nil {
		return nil, err
	}

	return &Runner{
		runtime: rtImpl,
		state:   state,
	}, nil
}

func (r *Runner) AlreadyDone(task model.InitTask) bool {
	if !task.Once {
		return false
	}
	return r.state.IsDone(task.Name)
}

func (r *Runner) Run(task model.InitTask) error {
	spec := runtime.RunSpec{
		Remove:   true,
		Volumes:  task.Volumes,
		Networks: task.Networks,
		Env:      task.Env,
		Image:    task.Image,
		Command:  task.Command,
	}

	hash := normalize.HashContainer(task.ToContainer())

	err := r.runtime.Run(spec, hash)
	if err != nil {
		return err
	}

	if task.Once {
		if err := r.state.MarkDone(task.Name); err != nil {
			return err
		}
	}

	return nil
}
