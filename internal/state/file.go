package state

import (
	"encoding/json"
	"os"
	"sync"
)

type FileStateStore struct {
	path  string
	state State
	mu    sync.Mutex
}

func NewFileStateStore(path string) (*FileStateStore, error) {
	s := &FileStateStore{
		path: path,
		state: State{
			InitTasks: make(map[string]bool),
		},
	}

	if err := s.Load(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *FileStateStore) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &s.state)
}

func (s *FileStateStore) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(s.state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.path, data, 0644)
}

func (s *FileStateStore) IsDone(name string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.state.InitTasks[name]
}

func (s *FileStateStore) MarkDone(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.state.InitTasks[name] = true
	return s.saveLocked()
}

func (s *FileStateStore) saveLocked() error {
	data, err := json.MarshalIndent(s.state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.path, data, 0644)
}
