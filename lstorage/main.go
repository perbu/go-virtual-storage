package lstorage

import (
	"io"
	"os"
)

type State struct {
}

func New() *State {
	s := &State{}
	return s
}

func (s *State) Create(name string) (io.WriteCloser, error) {
	return os.Create(name)
}
func (s *State) Open(name string) (io.ReadCloser, error) {
	return os.Open(name)
}
func (s *State) Remove(name string) error {
	return os.Remove(name)
}
func (s *State) Rename(oldName, newName string) error {
	panic("not implemented")
	return nil
}
func (s *State) List() ([]string, error) {
	panic("not implemented")
	return nil, nil
}
func (s *State) ListDir(dir string) ([]string, error) {
	panic("not implemented")
	return nil, nil
}
