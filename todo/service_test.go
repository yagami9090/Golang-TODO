package todo

import (
	"errors"
	"testing"
	"todo/entities"
)

type FakeRepoAdd struct {
	Repository
}

func (FakeRepoAdd) NewTask(*entities.Task) error {
	return errors.New("test")
}

func TestServiceError(t *testing.T) {
	srv := Service{FakeRepoAdd{}}
	err := srv.Add(entities.Task{})

	if err.Error() != "test" {
		t.Error("we expected error test")
	}
}
