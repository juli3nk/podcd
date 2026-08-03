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

func New(rt runtime.Backend, storageDir string, state state.Store) (InitRunner, error) {
	rtImpl, err := runtime.New(rt, storageDir)
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
	return r.state.IsDone(normalize.HashInitTask(task))
}

func (r *Runner) Run(task model.InitTask) error {
	hash := normalize.HashInitTask(task)

	spec := task.ToContainer()

	if err := r.runtime.Run(spec, hash); err != nil {
		return err
	}

	if task.Once {
		return r.state.MarkDone(hash)
	}

	return nil
}
