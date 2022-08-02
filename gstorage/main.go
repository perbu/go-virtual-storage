package gstorage

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"os"
)

const (
	bucketName = "perbu-cwy-test0"
	projectId  = "nimbus-testing-324411"
)

type State struct {
	client     *storage.Client
	projectID  string
	bucketName string
}

func NewClient(ctx context.Context) (*State, error) {
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "gcs-auth.json")
	if err != nil {
		return nil, fmt.Errorf("env error: %s", err)
	}
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	s := &State{
		client:     client,
		projectID:  projectId,
		bucketName: bucketName,
	}
	return s, nil
}

func (s *State) Create(name string) (io.WriteCloser, error) {
	wc := s.client.Bucket(s.bucketName).Object(name).NewWriter(context.TODO())
	return wc, nil
}
func (s *State) Open(name string) (io.ReadCloser, error) {
	rc, err := s.client.Bucket(s.bucketName).Object(name).NewReader(context.TODO())
	return rc, err
}
func (s *State) Remove(name string) error {
	err := s.client.Bucket(s.bucketName).Object(name).Delete(context.TODO())
	return err
}

func (s *State) Rename(oldName, newName string) error {
	panic("not implemented")
	return nil
}
func (s *State) List() ([]string, error) {
	var err error
	iterator := s.client.Bucket(s.bucketName).Objects(context.TODO(), nil)
	contents := make([]string, 0)
	for {
		oa, err := iterator.Next()
		if err != nil {
			break
		}
		contents = append(contents, oa.Name)
	}
	return contents, err
}

func (s *State) ListDir(dir string) ([]string, error) {
	panic("not implemented")
	return nil, nil
}
